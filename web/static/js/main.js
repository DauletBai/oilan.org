// oilan/web/static/js/main.js
document.addEventListener('DOMContentLoaded', () => {
    // We only need to run JS on the chat page.
    if (document.querySelector('#chat-section')) {
        handleChatPage();
    }
});

/**
 * Handles all logic for the main chat page.
 */
function handleChatPage() {
    let socket = null;

    const chatWindowBody = document.querySelector('#chat-window .card-body');
    const messageInput = document.getElementById('message-input');
    const sendButton = document.getElementById('send-button');
    const newChatButton = document.getElementById('new-chat-button');
    
    function addMessageToWindow(role, content) {
        const alignClass = role === 'user' ? 'text-end' : 'text-start';
        const colorClass = role === 'user' ? 'bg-primary text-white' : 'bg-secondary text-white';
        const messageWrapper = document.createElement('div');
        messageWrapper.className = `p-2 my-1 d-flex flex-column ${role === 'user' ? 'align-items-end' : 'align-items-start'}`;
        const messageDiv = document.createElement('div');
        messageDiv.className = `px-3 py-2 rounded-3`;
        messageDiv.style.maxWidth = '75%';
        messageDiv.classList.add(colorClass);
        messageDiv.textContent = content;
        messageWrapper.appendChild(messageDiv);
        chatWindowBody.appendChild(messageWrapper);
        chatWindowBody.scrollTop = chatWindowBody.scrollHeight;
    }

    function connectWebSocket() {
        if (socket) { socket.close(); }
        
        const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        // We no longer need to pass any token here. The browser sends the secure cookie automatically.
        const wsURL = `${proto}//${window.location.host}/ws/chat`;

        socket = new WebSocket(wsURL);

        socket.onopen = () => {
            console.log('WebSocket connection established.');
            messageInput.disabled = false;
            sendButton.disabled = false;
            newChatButton.disabled = false;
            messageInput.focus();
        };

        socket.onmessage = (event) => {
            addMessageToWindow('ai', event.data);
            sendButton.disabled = false;
            messageInput.disabled = false;
            messageInput.focus();
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed.');
            addMessageToWindow('ai', 'Connection has been closed.');
            messageInput.disabled = true;
            sendButton.disabled = true;
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            addMessageToWindow('ai', 'A connection error occurred.');
        };
    }

    function sendMessage() {
        const content = messageInput.value.trim();
        if (!content || !socket || socket.readyState !== WebSocket.OPEN) return;

        addMessageToWindow('user', content);
        socket.send(content);
        
        messageInput.value = '';
        messageInput.disabled = true;
        sendButton.disabled = true;
    }

    sendButton.addEventListener('click', sendMessage);
    messageInput.addEventListener('keyup', (event) => {
        if (event.key === 'Enter') { sendMessage(); }
    });
    newChatButton.addEventListener('click', () => {
        chatWindowBody.innerHTML = '';
        connectWebSocket();
    });

    // --- Initial Load ---
    connectWebSocket();
}