'use client';

import { useState, useEffect } from 'react';

interface Conversation {
    id: string;
    contact: {
        id: string;
        name: string;
        avatar_url?: string;
    };
    platform: string;
    last_message_text: string;
    last_message_at: string;
    unread_count: number;
}

interface Message {
    id: string;
    content: string;
    direction: string;
    created_at: string;
}

export default function InboxPage() {
    const [conversations, setConversations] = useState<Conversation[]>([]);
    const [selectedConv, setSelectedConv] = useState<string | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchConversations();
    }, []);

    const fetchConversations = async () => {
        try {
            const res = await fetch('/api/conversations');
            const data = await res.json();
            setConversations(data.conversations || []);
        } catch (error) {
            console.error('Failed to fetch conversations:', error);
        } finally {
            setLoading(false);
        }
    };

    const fetchMessages = async (convId: string) => {
        try {
            const res = await fetch(`/api/conversations/${convId}`);
            const data = await res.json();
            setMessages(data.messages || []);
        } catch (error) {
            console.error('Failed to fetch messages:', error);
        }
    };

    const selectConversation = (id: string) => {
        setSelectedConv(id);
        fetchMessages(id);
    };

    const sendMessage = async () => {
        if (!newMessage.trim() || !selectedConv) return;

        const conv = conversations.find(c => c.id === selectedConv);
        if (!conv) return;

        try {
            await fetch('/api/messages', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    conversation_id: selectedConv,
                    platform: conv.platform,
                    recipient_id: conv.contact.id,
                    content: newMessage,
                    content_type: 'text',
                }),
            });

            setNewMessage('');
            fetchMessages(selectedConv);
        } catch (error) {
            console.error('Failed to send message:', error);
        }
    };

    const getInitials = (name: string) => {
        return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
    };

    const formatTime = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
    };

    if (loading) {
        return (
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
                <div className="loading-spinner"></div>
            </div>
        );
    }

    return (
        <div style={{ display: 'flex', height: '100vh' }}>
            {/* Conversation List */}
            <div style={{
                width: '350px',
                borderRight: '1px solid var(--color-border)',
                display: 'flex',
                flexDirection: 'column'
            }}>
                <div style={{ padding: 'var(--space-md)', borderBottom: '1px solid var(--color-divider)' }}>
                    <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, marginBottom: 'var(--space-md)' }}>
                        Inbox
                    </h2>
                    <div className="search-container">
                        <span className="search-icon">🔍</span>
                        <input
                            type="text"
                            className="search-input"
                            placeholder="Search conversations..."
                        />
                    </div>
                </div>

                <div className="conversation-list" style={{ flex: 1, overflowY: 'auto' }}>
                    {conversations.length === 0 ? (
                        <div className="empty-state" style={{ padding: 'var(--space-xl)' }}>
                            <div className="empty-state-icon">💬</div>
                            <div className="empty-state-title">No conversations yet</div>
                            <div className="empty-state-description">
                                Incoming messages will appear here
                            </div>
                        </div>
                    ) : (
                        conversations.map((conv) => (
                            <div
                                key={conv.id}
                                className={`conversation-item ${selectedConv === conv.id ? 'active' : ''}`}
                                onClick={() => selectConversation(conv.id)}
                            >
                                <div className="conversation-avatar">
                                    {getInitials(conv.contact.name)}
                                </div>
                                <div className="conversation-content">
                                    <div className="conversation-header">
                                        <span className="conversation-name">{conv.contact.name}</span>
                                        <span className="conversation-time">{formatTime(conv.last_message_at)}</span>
                                    </div>
                                    <div className="conversation-preview">{conv.last_message_text}</div>
                                </div>
                                <div className="conversation-meta">
                                    <span className={`platform-badge ${conv.platform}`}>
                                        {conv.platform === 'whatsapp' ? 'WA' : 'IG'}
                                    </span>
                                    {conv.unread_count > 0 && (
                                        <span className="unread-badge">{conv.unread_count}</span>
                                    )}
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>

            {/* Chat View */}
            <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                {selectedConv ? (
                    <>
                        <div className="chat-header">
                            <div className="conversation-avatar">
                                {getInitials(conversations.find(c => c.id === selectedConv)?.contact.name || '')}
                            </div>
                            <div>
                                <div style={{ fontWeight: 600 }}>
                                    {conversations.find(c => c.id === selectedConv)?.contact.name}
                                </div>
                                <div style={{ fontSize: 'var(--text-xs)', color: 'var(--color-text-secondary)' }}>
                                    {conversations.find(c => c.id === selectedConv)?.platform}
                                </div>
                            </div>
                        </div>

                        <div className="chat-messages">
                            {messages.map((msg) => (
                                <div key={msg.id} className={`message-bubble ${msg.direction}`}>
                                    <div className="message-text">{msg.content}</div>
                                    <div className="message-time">{formatTime(msg.created_at)}</div>
                                </div>
                            ))}
                        </div>

                        <div className="chat-input-container">
                            <div className="chat-input-wrapper">
                                <input
                                    className="chat-input"
                                    placeholder="Type a message..."
                                    value={newMessage}
                                    onChange={(e) => setNewMessage(e.target.value)}
                                    onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                                />
                                <button className="send-button" onClick={sendMessage}>
                                    ➤
                                </button>
                            </div>
                        </div>
                    </>
                ) : (
                    <div className="empty-state" style={{ flex: 1 }}>
                        <div className="empty-state-icon">💬</div>
                        <div className="empty-state-title">Select a conversation</div>
                        <div className="empty-state-description">
                            Choose a conversation from the list to start chatting
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
