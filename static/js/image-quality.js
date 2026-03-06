/**
 * image-quality.js — Smart Image Validation & Compression
 *
 * All processing runs on-device via Canvas API (no server round-trip).
 * Designed for fast execution on mid-range Android devices.
 *
 * Public API (window.ImageQuality):
 *   analyze(blob)                        → Promise<QualityResult>
 *   compress(blob, maxWidth?, quality?)  → Promise<Blob>
 *   showFeedback(containerEl, result)    → void
 *   clearFeedback(containerEl)           → void
 *
 * QualityResult: { sharpness, brightness, sizeKB, warnings[], passed }
 */
(function () {
  'use strict';

  /* ── Thresholds ──────────────────────────────────────────────────────────── */
  var ANALYSIS_MAX_W   = 256;   // downsample to this width before analysis (speed)
  var COMPRESS_MAX_W   = 2000;  // max output width after compression
  var COMPRESS_QUALITY = 0.82;  // JPEG encode quality (0–1)
  var MIN_SIZE_KB      = 20;    // file too small → likely a thumbnail
  var MAX_SIZE_MB      = 8;     // above this we always compress
  var MIN_SHARPNESS    = 100;   // Laplacian variance threshold
  var MIN_BRIGHTNESS   = 40;    // average luminance (0–255)
  var MAX_BRIGHTNESS   = 225;   // average luminance (0–255)

  /* ── Async image decode ──────────────────────────────────────────────────── */
  function _decode(blob) {
    if (typeof createImageBitmap !== 'undefined') {
      return createImageBitmap(blob);
    }
    return new Promise(function (resolve, reject) {
      var img = new Image();
      var url = URL.createObjectURL(blob);
      img.onload  = function () { URL.revokeObjectURL(url); resolve(img); };
      img.onerror = function () { URL.revokeObjectURL(url); reject(new Error('decode failed')); };
      img.src = url;
    });
  }

  /* ── Draw to a small analysis canvas ────────────────────────────────────── */
  function _toSmallCanvas(bmp) {
    var srcW  = bmp.naturalWidth  || bmp.width;
    var srcH  = bmp.naturalHeight || bmp.height;
    var scale = Math.min(1, ANALYSIS_MAX_W / srcW);
    var w     = Math.max(1, Math.round(srcW * scale));
    var h     = Math.max(1, Math.round(srcH * scale));
    var c     = document.createElement('canvas');
    c.width   = w;
    c.height  = h;
    c.getContext('2d', { willReadFrequently: true }).drawImage(bmp, 0, 0, w, h);
    return c;
  }

  /* ── Convert to grayscale Float32 array ─────────────────────────────────── */
  function _grayscale(ctx, w, h) {
    var d    = ctx.getImageData(0, 0, w, h).data;
    var gray = new Float32Array(w * h);
    for (var i = 0; i < w * h; i++) {
      var p = i * 4;
      gray[i] = d[p] * 0.299 + d[p + 1] * 0.587 + d[p + 2] * 0.114;
    }
    return gray;
  }

  /* ── Laplacian variance (sharpness score) ────────────────────────────────── */
  function _laplacianVariance(gray, w, h) {
    var laps = [];
    for (var y = 1; y < h - 1; y++) {
      for (var x = 1; x < w - 1; x++) {
        var c  = gray[ y      * w + x    ];
        var tl = gray[(y - 1) * w + (x-1)];
        var t  = gray[(y - 1) * w + x    ];
        var tr = gray[(y - 1) * w + (x+1)];
        var l  = gray[ y      * w + (x-1)];
        var r  = gray[ y      * w + (x+1)];
        var bl = gray[(y + 1) * w + (x-1)];
        var b  = gray[(y + 1) * w + x    ];
        var br = gray[(y + 1) * w + (x+1)];
        laps.push(8*c - tl - t - tr - l - r - bl - b - br);
      }
    }
    var n    = laps.length;
    if (n === 0) return 0;
    var mean = 0;
    for (var i = 0; i < n; i++) mean += laps[i];
    mean /= n;
    var variance = 0;
    for (var i = 0; i < n; i++) variance += (laps[i] - mean) * (laps[i] - mean);
    return variance / n;
  }

  /* ── Average brightness ──────────────────────────────────────────────────── */
  function _avgBrightness(gray) {
    var sum = 0;
    for (var i = 0; i < gray.length; i++) sum += gray[i];
    return gray.length > 0 ? sum / gray.length : 128;
  }

  /* ════════════════════════════════════════════════════════════════════════════
   * PUBLIC: analyzeImageQuality(blob)
   * Returns Promise<{ sharpness, brightness, sizeKB, warnings[], passed }>
   * ════════════════════════════════════════════════════════════════════════════ */
  function analyzeImageQuality(blob) {
    var sizeKB   = blob.size / 1024;
    var warnings = [];
    var passed   = true;

    if (blob.type === 'application/pdf') {
      return Promise.resolve({ sharpness: 999, brightness: 128, sizeKB: sizeKB, warnings: [], passed: true, isPDF: true });
    }

    return _decode(blob).then(function (bmp) {
      var canvas     = _toSmallCanvas(bmp);
      var ctx        = canvas.getContext('2d', { willReadFrequently: true });
      var gray       = _grayscale(ctx, canvas.width, canvas.height);
      var sharpness  = _laplacianVariance(gray, canvas.width, canvas.height);
      var brightness = _avgBrightness(gray);

      if (bmp.close) bmp.close();

      if (sharpness < MIN_SHARPNESS) {
        warnings.push('ภาพเบลอเกินไป กรุณาถ่ายใหม่ให้ชัดขึ้น');
        passed = false;
      }
      if (brightness < MIN_BRIGHTNESS) {
        warnings.push('ภาพมืดเกินไป กรุณาเพิ่มแสงหรือถ่ายในที่สว่างกว่านี้');
        passed = false;
      }
      if (brightness > MAX_BRIGHTNESS) {
        warnings.push('ภาพสว่างจนเกินไป (Overexposed) กรุณาหลีกเลี่ยงแสงจ้าตรงหน้า');
        passed = false;
      }
      if (sizeKB < MIN_SIZE_KB) {
        warnings.push('ไฟล์มีขนาดเล็กเกินไป อาจขาดรายละเอียดที่จำเป็น');
        passed = false;
      }

      return { sharpness: sharpness, brightness: brightness, sizeKB: sizeKB, warnings: warnings, passed: passed };

    }).catch(function (e) {
      console.warn('[ImageQuality] analyze error:', e);
      return { sharpness: 0, brightness: 128, sizeKB: sizeKB, warnings: [], passed: true };
    });
  }

  /* ════════════════════════════════════════════════════════════════════════════
   * PUBLIC: compressImage(blob, maxWidth?, quality?)
   * Returns Promise<Blob> — resized + JPEG compressed.
   * PDFs pass through unchanged.
   * ════════════════════════════════════════════════════════════════════════════ */
  function compressImage(blob, maxWidth, quality) {
    maxWidth = maxWidth || COMPRESS_MAX_W;
    quality  = (quality !== undefined) ? quality : COMPRESS_QUALITY;

    if (blob.type === 'application/pdf') return Promise.resolve(blob);

    return _decode(blob).then(function (bmp) {
      var srcW  = bmp.naturalWidth  || bmp.width;
      var srcH  = bmp.naturalHeight || bmp.height;
      var scale = Math.min(1, maxWidth / srcW);
      var dstW  = Math.max(1, Math.round(srcW * scale));
      var dstH  = Math.max(1, Math.round(srcH * scale));

      var canvas   = document.createElement('canvas');
      canvas.width  = dstW;
      canvas.height = dstH;
      canvas.getContext('2d').drawImage(bmp, 0, 0, dstW, dstH);
      if (bmp.close) bmp.close();

      return new Promise(function (resolve) {
        canvas.toBlob(function (b) { resolve(b || blob); }, 'image/jpeg', quality);
      });
    }).catch(function (e) {
      console.warn('[ImageQuality] compress error:', e);
      return blob;
    });
  }

  /* ════════════════════════════════════════════════════════════════════════════
   * PUBLIC: showQualityFeedback(containerEl, result)
   * Renders a quality card into the given container element.
   * ════════════════════════════════════════════════════════════════════════════ */
  function showQualityFeedback(container, result) {
    if (!container) return;
    if (result.isPDF) { container.innerHTML = ''; return; }

    var ok     = result.passed;
    var accent = ok ? '#059669' : '#d97706';
    var bg     = ok ? '#f0fdf4' : '#fffbeb';
    var border = ok ? '#bbf7d0' : '#fde68a';
    var icon   = ok ? '✅' : '⚠️';
    var title  = ok
      ? 'ภาพผ่านการตรวจสอบ'
      : 'ภาพเบลอ/มืดเกินไป กรุณาถ่ายใหม่เพื่อลดการถูกตีกลับงาน';

    var sharpPct  = Math.min(100, (result.sharpness  / 500) * 100).toFixed(0);
    var brightPct = Math.min(100, (result.brightness / 255) * 100).toFixed(0);
    var sharpColor  = result.sharpness  < MIN_SHARPNESS                                  ? '#ef4444' : '#10b981';
    var brightColor = (result.brightness < MIN_BRIGHTNESS || result.brightness > MAX_BRIGHTNESS) ? '#f59e0b' : '#10b981';
    var sharpTxt  = result.sharpness  < MIN_SHARPNESS  ? '#dc2626' : '#374151';
    var brightTxt = (result.brightness < MIN_BRIGHTNESS || result.brightness > MAX_BRIGHTNESS) ? '#dc2626' : '#374151';

    var sizeFmt = result.sizeKB >= 1024
      ? (result.sizeKB / 1024).toFixed(1) + ' MB'
      : Math.round(result.sizeKB) + ' KB';
    var sizeTxt = result.sizeKB < MIN_SIZE_KB ? '#dc2626' : '#374151';

    var warningsHtml = '';
    if (result.warnings.length > 0) {
      warningsHtml = '<div style="margin-top:6px;color:#b45309;font-size:0.82rem;line-height:1.5;">'
        + result.warnings.map(function (w) { return '• ' + w; }).join('<br>')
        + '</div>';
    }

    container.innerHTML = '<style>'
      + '.iq-bar-track{background:#e5e7eb;border-radius:99px;height:6px;overflow:hidden;margin-top:2px;}'
      + '.iq-bar-fill{height:100%;border-radius:99px;transition:width .5s ease;}'
      + '@keyframes iq-in{from{opacity:0;transform:translateY(-4px)}to{opacity:1;transform:translateY(0)}}'
      + '</style>'
      + '<div style="border:1.5px solid ' + border + ';border-radius:12px;background:' + bg + ';padding:12px 14px;font-family:inherit;font-size:0.88rem;animation:iq-in .3s ease;">'
      +   '<div style="font-weight:700;color:' + accent + ';margin-bottom:8px;">' + icon + ' ' + title + '</div>'
      +   '<div style="display:grid;gap:6px;">'
      +     '<div>'
      +       '<div style="display:flex;justify-content:space-between;">'
      +         '<span>🔍 ความชัด</span>'
      +         '<span style="color:' + sharpTxt + ';">' + Math.round(result.sharpness) + '</span>'
      +       '</div>'
      +       '<div class="iq-bar-track"><div class="iq-bar-fill" style="width:' + sharpPct + '%;background:' + sharpColor + ';"></div></div>'
      +     '</div>'
      +     '<div>'
      +       '<div style="display:flex;justify-content:space-between;">'
      +         '<span>☀️ ความสว่าง</span>'
      +         '<span style="color:' + brightTxt + ';">' + Math.round(result.brightness) + '/255</span>'
      +       '</div>'
      +       '<div class="iq-bar-track"><div class="iq-bar-fill" style="width:' + brightPct + '%;background:' + brightColor + ';"></div></div>'
      +     '</div>'
      +     '<div style="display:flex;justify-content:space-between;">'
      +       '<span>💾 ขนาดไฟล์</span>'
      +       '<span style="color:' + sizeTxt + ';">' + sizeFmt + '</span>'
      +     '</div>'
      +     warningsHtml
      +   '</div>'
      + '</div>';

    container.style.display = 'block';
  }

  /* ════════════════════════════════════════════════════════════════════════════
   * PUBLIC: clearQualityFeedback(containerEl)
   * ════════════════════════════════════════════════════════════════════════════ */
  function clearQualityFeedback(container) {
    if (!container) return;
    container.innerHTML = '';
    container.style.display = 'none';
  }

  /* ── Expose ────────────────────────────────────────────────────────────────── */
  window.ImageQuality = {
    analyze:       analyzeImageQuality,
    compress:      compressImage,
    showFeedback:  showQualityFeedback,
    clearFeedback: clearQualityFeedback
  };

}());
