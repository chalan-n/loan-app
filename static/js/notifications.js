// Notification System Logic
let ws;
let notifHistory = [];

function connectWebSocket() {
    const wsProtocol = location.protocol === 'https:' ? 'wss://' : 'ws://';
    ws = new WebSocket(`${wsProtocol}${location.host}/ws`);

    ws.onopen = () => console.log('WebSocket เชื่อมต่อสำเร็จ');

    ws.onmessage = function (event) {
        const data = JSON.parse(event.data);
        addNotification(data.title || 'แจ้งเตือนใหม่!', data.message || '');

        new Audio('/static/sound/notif.mp3').play().catch(() => { });
    };

    ws.onclose = () => {
        console.log('WebSocket หลุด! เชื่อมใหม่ใน 3 วินาที...');
        setTimeout(connectWebSocket, 3000);
    };
}

document.addEventListener("DOMContentLoaded", function() {
    connectWebSocket();

    const notifBell = document.getElementById('notificationBell');
    const notifDropdown = document.getElementById('notifDropdown');
    const notifBadge = document.getElementById('notifBadge');
    
    // คลิกกระดิ่ง
    if(notifBell) {
        notifBell.onclick = (e) => {
            e.stopPropagation();
            toggleNotifDropdown();
        };
    }
    
    if(notifDropdown) {
        notifDropdown.onclick = (e) => e.stopPropagation();
    }
    
    document.addEventListener('click', () => {
        if(notifDropdown) notifDropdown.style.display = 'none';
    });
});

function toggleNotifDropdown() {
    const dropdown = document.getElementById('notifDropdown');
    if (!dropdown) return;
    
    if (dropdown.style.display === 'block') {
        dropdown.style.display = 'none';
    } else {
        dropdown.style.display = 'block';
        // Reset unread status
        notifHistory.forEach(n => n.unread = false);
        updateNotifList();
        updateBadge();
        updateClearButtonState();
    }
}

function clearAllNotifications() {
    if(!confirm('ล้างประวัติทั้งหมดหรือไม่?')) return;
    
    notifHistory = [];
    updateNotifList();
    updateBadge();
    updateClearButtonState();
}

function updateClearButtonState() {
    const btn = document.getElementById('clearAllNotif');
    if(!btn) return;
    
    if (notifHistory.length > 0) {
        btn.disabled = false;
        btn.style.opacity = '1';
        btn.style.cursor = 'pointer';
    } else {
        btn.disabled = true;
        btn.style.opacity = '0.5';
        btn.style.cursor = 'not-allowed';
    }
}

function addNotification(title, message) {
    const now = new Date();
    const timeStr = now.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' });

    const notif = { title, message, time: timeStr, timestamp: now.getTime(), unread: true };
    notifHistory.unshift(notif);
    if (notifHistory.length > 50) notifHistory.pop();

    updateNotifList();
    updateBadge();
    updateClearButtonState();

    // แสดง Toast
    const toastTitle = document.getElementById('toastTitle');
    const toastMsg = document.getElementById('toastMsg');
    const toast = document.getElementById('realtimeToast');
    
    if(toast && toastTitle && toastMsg) {
        toastTitle.textContent = title;
        toastMsg.textContent = message;
        toast.classList.add('show');
        if(toast.hideTimer) clearTimeout(toast.hideTimer);
        toast.hideTimer = setTimeout(() => toast.classList.remove('show'), 5000);
    }
}

function updateNotifList() {
    const list = document.getElementById('notifList');
    if(!list) return;

    if (notifHistory.length === 0) {
        list.innerHTML = '<div style="text-align:center;color:#9ca3af;padding:40px;font-size:0.95rem;" id="hasNoNotif">ยังไม่มีแจ้งเตือน</div>';
        return;
    }
    list.innerHTML = notifHistory.map(n => `
      <div class="notif-item ${n.unread ? 'unread' : ''}" onclick="markAsRead(this, ${n.timestamp})">
        <div class="title">${n.title}</div>
        <div class="message">${n.message}</div>
        <div class="time">${n.time}</div>
      </div>
    `).join('');
}

function updateBadge() {
    const badge = document.getElementById('notifBadge');
    if(!badge) return;

    const unread = notifHistory.filter(n => n.unread).length;
    badge.textContent = unread;
    badge.style.display = unread > 0 ? 'flex' : 'none';
}

function markAsRead(el, ts) {
    const n = notifHistory.find(x => x.timestamp === ts);
    if (n) n.unread = false;
    el.classList.remove('unread');
    updateBadge();
}
