// Custom Loading Overlay Logic

document.addEventListener("DOMContentLoaded", function () {
    // Create overlay HTML if not exists
    if (!document.getElementById('loadingOverlay')) {
        const overlay = document.createElement('div');
        overlay.id = 'loadingOverlay';
        overlay.innerHTML = `
      <div class="loading-spinner"></div>
      <div class="loading-text">กำลังประมวลผล...</div>
    `;
        document.body.appendChild(overlay);
    }
});

function showLoading(text = 'กำลังประมวลผล...') {
    const overlay = document.getElementById('loadingOverlay');
    if (overlay) {
        overlay.querySelector('.loading-text').textContent = text;
        overlay.classList.add('show');
    }
}

function hideLoading() {
    const overlay = document.getElementById('loadingOverlay');
    if (overlay) {
        overlay.classList.remove('show');
    }
}

// Auto-attach to forms
document.querySelectorAll('form').forEach(form => {
    form.addEventListener('submit', function (e) {
        // If client-side validation fails (and prevents default), we shouldn't show loading.
        // However, if the form logic prevents default itself (custom validation), we need to manually call showLoading there.
        // This generic listener assumes standard submit.
        // If the form is submitted via JS manually, this might not trigger or we need to call showLoading manually.

        // Check if the form has 'onsubmit' attribute that might return false
        // But basic HTML5 validation or custom JS usually handles valid check.
        // We'll set a small timeout to checking if submission was prevented.

        // For our specific use case in Step 1-7, validation prevents default if invalid.
        // So if we reach here, we are submitting (unless preventDefault called immediately).

        // Let's delay slightly to allow other validation listeners to run.
        // Actually, closest listener runs first in capturing, but we are bubbling.
        // If invalid, e.preventDefault() was called by the validation script.

        // Wait for current event loop to finish to see if defaultPrevented
        /*
        setTimeout(() => {
          if (!e.defaultPrevented) {
            showLoading('กำลังบันทึกข้อมูล...');
          }
        }, 0);
        */
        // The above generic approach can be tricky. 
        // It is safer to rely on explicit calls or assume if this listener runs, we show it, 
        // UNLESS the validation script stops propagation.
        // But our validation script in step1 uses e.preventDefault() if invalid.
        // If we bind this listener *after* the validation listener, we can check e.defaultPrevented.
    });
});
