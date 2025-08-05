// oilan/web/static/js/main.js
document.addEventListener('DOMContentLoaded', () => {
    // Simple router: call the correct handler based on which page is loaded.
    if (document.querySelector('#auth-section')) {
        // We are on the welcome page, no JS needed for now.
    } else if (document.querySelector('#chat-section')) {
        handleChatPage();
    }
});

/**
 * Handles logic for the welcome/login page.
 */
function handleWelcomePage() {
    // This part handles the redirect from Google OAuth.
    // After a successful login, the backend redirects here and provides the JWT.
    if (window.location.pathname.includes('/auth/google/callback')) {
        const pageContent = document.body.innerText;
        try {
            const tokenData = JSON.parse(pageContent);
            if (tokenData.token) {
                // If we find a token, we save it and redirect to the chat page.
                localStorage.setItem('oilan_jwt', tokenData.token);
                window.location.href = '/chat';
                return;
            }
        } catch (e) {
            // If the content is not a valid JSON, it's likely an error message.
            document.body.innerHTML = `<h1>Login Failed</h1><p>An error occurred during authentication.</p><a href='/'>Try again</a>`;
        }
    }

    // If the user already has a token, they shouldn't see the login page.
    const jwtToken = localStorage.getItem('oilan_jwt');
    if (jwtToken) {
        window.location.href = '/chat';
    }
}

/**
 * Handles all logic for the main chat page.
 */
function handleChatPage() {
    // State Management
    let currentDialogID = null;
    let socket = null;

    // DOM Elements
    const chatWindowBody = document.querySelector('#chat-window .card-body');
    const messageInput = document.getElementById('message-input');
    const sendButton = document.getElementById('send-button');
    const newChatButton = document.getElementById('new-chat-button');

    /**
     * Appends a message to the chat window UI.
     * @param {string} role - 'user' or 'ai'.
     * @param {string} content - The text of the message.
     */
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

    /**
     * A helper function for making authenticated API calls.
     * @param {string} endpoint - The API endpoint (e.g., '/api/v1/dialogs').
     * @param {string} method - The HTTP method (e.g., 'POST').
     * @param {object} body - The JSON body for the request.
     */
    async function apiFetch(endpoint, method, body) {
        const headers = { 'Content-Type': 'application/json' };
        const response = await fetch(endpoint, {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined,
        });
        if (!response.ok) {
            // If we get an unauthorized error, redirect to login page.
            if (response.status === 401) {
                window.location.href = '/';
                return;
            }
            const error = await response.json();
            throw new Error(error.error);
        }
        return response.json();
    }
    
    /**
     * Creates a new dialog session and establishes a WebSocket connection.
     */
    async function startNewChat() {
        chatWindowBody.innerHTML = '';
        addMessageToWindow('ai', 'Creating a new secure session...');
        messageInput.disabled = true;
        sendButton.disabled = true;
        try {
            const dialog = await apiFetch('/api/v1/dialogs', 'POST', { title: 'New Web Chat' });
            currentDialogID = dialog.id;
            // Now connect to WebSocket
            connectWebSocket();
        } catch (error) {
            addMessageToWindow('ai', `Error: Could not start a new chat. ${error.message}`);
        }
    }

    /**
     * Connects to the WebSocket server.
     */
    function connectWebSocket() {
        if (socket) {
            socket.close();
        }
        
        // Construct WebSocket URL. Use wss:// for secure connections (in production).
        const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsURL = `${proto}//${window.location.host}/ws/chat?dialogID=${currentDialogID}`;

        socket = new WebSocket(wsURL);

        socket.onopen = () => {
            console.log('WebSocket connection established.');
            chatWindowBody.innerHTML = ''; // Clear "creating session" message
            addMessageToWindow('ai', 'Hello! I am ready. How can I help you today?');
            messageInput.disabled = false;
            sendButton.disabled = false;
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

    // --- Event Listeners ---
    sendButton.addEventListener('click', sendMessage);
    messageInput.addEventListener('keyup', (event) => {
        if (event.key === 'Enter') {
            sendMessage();
        }
    });
    newChatButton.addEventListener('click', startNewChat);

    // --- Initial Load ---
    startNewChat();
}
