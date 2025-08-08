// oilan/web/static/js/main.js
document.addEventListener('DOMContentLoaded', () => {
    if (document.querySelector('#chat-section')) {
        handleChatPage();
    } else {
        // This is for the login page
        handleWelcomePage();
    }
});

function handleWelcomePage() {
    fetch('/api/v1/session').then(response => {
        if (response.ok) { window.location.href = '/chat'; }
    });
}

async function handleChatPage() {
    let currentDialogID = null;
    let socket = null;

    const chatWindowCard = document.getElementById('chat-window');
    const chatWindowBody = document.querySelector('#chat-window .card-body');
    const messageInput = document.getElementById('message-input');
    const sendButton = document.getElementById('send-button');
    const newChatButton = document.getElementById('new-chat-button');
    const dialogList = document.getElementById('dialog-list');

    /**
     * Appends a message to the chat window UI.
     */
    function addMessageToWindow(role, content) {
        const messageWrapper = document.createElement('div');
        messageWrapper.className = `p-2 my-1 d-flex flex-column ${role === 'user' ? 'align-items-end' : 'align-items-start'}`;

        const messageDiv = document.createElement('div');
        messageDiv.className = 'px-3 py-2 rounded-3';
        messageDiv.style.maxWidth = '75%';
        // messageDiv.classList.add(...colorClass.split(' '));

        // We add classes one by one, not as a single string.
        if (role === 'user') {
            messageDiv.classList.add('bg-primary', 'text-white');
        } else {
            messageDiv.classList.add('bg-light', 'text-dark', 'border');
        }
        messageDiv.textContent = content;
        messageWrapper.appendChild(messageDiv);
        chatWindowBody.appendChild(messageWrapper);
        
        // chatWindowBody.scrollTop = chatWindowBody.scrollHeight;
        requestAnimationFrame(() => {
            chatWindowCard.scrollTo({ top: chatWindowCard.scrollHeight, behavior: 'smooth' });
        });
    }

    /**
     * Helper for making authenticated API calls.
     */
    async function apiFetch(endpoint, method, body) {
        const headers = { 'Content-Type': 'application/json' };
        const response = await fetch('/api/v1' + endpoint, {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined,
        });
        if (!response.ok) {
            if (response.status === 401) { window.location.href = '/'; }
            const error = await response.json();
            throw new Error(error.error);
        }
        if (response.status === 204) return null;
        return response.json();
    }
    
    /**
     * Connects to the WebSocket server.
     */
    function connectWebSocket(dialogID) {
        if (socket) { socket.close(); }
        currentDialogID = dialogID;

        const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsURL = `${proto}//${window.location.host}/ws/chat?dialogID=${dialogID}`;

        socket = new WebSocket(wsURL);

        socket.onopen = () => {
            console.log('WebSocket connection established for dialog', dialogID);
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

    /**
     * Loads a specific dialog's history and connects to it.
     */
    async function loadDialog(dialogID) {
        chatWindowBody.innerHTML = '';
        addMessageToWindow('ai', 'Loading history...');
        try {
            const dialogData = await apiFetch(`/dialogs/${dialogID}`, 'GET');
            chatWindowBody.innerHTML = ''; // Clear loading message
            if (dialogData.messages && dialogData.messages.length > 0) {
                 dialogData.messages.forEach(msg => addMessageToWindow(msg.role, msg.content));
            } else {
                addMessageToWindow('ai', 'This is a new chat. How can I help?');
            }
            connectWebSocket(dialogID);
        } catch (error) {
            addMessageToWindow('ai', `Error loading chat: ${error.message}`);
        }
    }

    /**
     * Creates a new dialog session and connects to it.
     */
    async function startNewChat() {
        chatWindowBody.innerHTML = ''; 
        addMessageToWindow('ai', 'Creating a new session...');
        messageInput.disabled = true;
        sendButton.disabled = true;
        try {
            const dialog = await apiFetch('/dialogs', 'POST', { title: 'New Chat' });
            await loadUserDialogs(); // Refresh dialog list
            await loadDialog(dialog.id); // Load the new (empty) dialog
        } catch (error) {
            addMessageToWindow('ai', `Error: ${error.message}`);
        }
    }

    /**
     * Sends a message over the WebSocket connection.
     */
    function sendMessage() {
        const content = messageInput.value.trim();
        if (!content || !socket || socket.readyState !== WebSocket.OPEN) return;
        addMessageToWindow('user', content);
        socket.send(content);
        messageInput.value = '';
        messageInput.disabled = true;
        sendButton.disabled = true;
    }

    /**
     * Renders the list of dialogs in the sidebar.
     */
    function renderDialogList(dialogs) {
        dialogList.innerHTML = '';
        if (dialogs && dialogs.length > 0) {
            dialogs.forEach(dialog => {
                const item = document.createElement('a');
                item.href = '#';
                item.className = 'list-group-item list-group-item-action';
                item.textContent = dialog.title || `Chat ${dialog.id}`;
                item.dataset.dialogId = dialog.id;
                if (dialog.id === currentDialogID) {
                    item.classList.add('active');
                }
                dialogList.appendChild(item);
            });
        } else {
            dialogList.innerHTML = '<p class="text-muted p-2">No past conversations.</p>';
        }
    }

    /**
     * Fetches all dialogs for the user and renders them.
     */
    async function loadUserDialogs() {
        try {
            const dialogs = await apiFetch('/dialogs', 'GET');
            renderDialogList(dialogs);
        } catch (error) {
            console.error("Failed to load dialogs:", error.message);
        }
    }

    // Event Listeners
    sendButton.addEventListener('click', sendMessage);
    messageInput.addEventListener('keyup', (event) => {
        if (event.key === 'Enter') { sendMessage(); }
    });
    newChatButton.addEventListener('click', startNewChat);
    
    dialogList.addEventListener('click', (event) => {
        event.preventDefault();
        const dialogId = event.target.closest('a')?.dataset.dialogId;
        if (dialogId) {
            loadDialog(parseInt(dialogId, 10));
        }
    });

    // Initial Load
    const initialDialogs = await apiFetch('/dialogs', 'GET');
    renderDialogList(initialDialogs);
    if (initialDialogs && initialDialogs.length > 0) {
        // Load the most recent dialog
        loadDialog(initialDialogs[0].id);
    } else {
        // Or start a new one if there's no history
        startNewChat();
    }
}