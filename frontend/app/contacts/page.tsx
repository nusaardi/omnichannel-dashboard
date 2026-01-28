'use client';

import { useState, useEffect } from 'react';

interface Contact {
    id: string;
    name: string;
    phone?: string;
    email?: string;
    whatsapp_id?: string;
    instagram_id?: string;
    created_at: string;
}

export default function ContactsPage() {
    const [contacts, setContacts] = useState<Contact[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');

    useEffect(() => {
        fetchContacts();
    }, []);

    const fetchContacts = async () => {
        try {
            const res = await fetch('/api/contacts');
            const data = await res.json();
            setContacts(data.contacts || []);
        } catch (error) {
            console.error('Failed to fetch contacts:', error);
        } finally {
            setLoading(false);
        }
    };

    const getInitials = (name: string) => {
        return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
    };

    const filteredContacts = contacts.filter(contact =>
        contact.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        contact.phone?.includes(searchQuery) ||
        contact.email?.toLowerCase().includes(searchQuery.toLowerCase())
    );

    if (loading) {
        return (
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
                <div className="loading-spinner"></div>
            </div>
        );
    }

    return (
        <div className="page-container">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--space-lg)' }}>
                <h1 style={{ fontSize: 'var(--text-2xl)', fontWeight: 700 }}>
                    Contacts
                </h1>
                <button className="btn btn-primary">
                    + Add Contact
                </button>
            </div>

            <div style={{ marginBottom: 'var(--space-lg)' }}>
                <div className="search-container" style={{ maxWidth: '400px' }}>
                    <span className="search-icon">🔍</span>
                    <input
                        type="text"
                        className="search-input"
                        placeholder="Search contacts..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                </div>
            </div>

            <div className="card">
                {filteredContacts.length === 0 ? (
                    <div className="empty-state">
                        <div className="empty-state-icon">👥</div>
                        <div className="empty-state-title">No contacts found</div>
                        <div className="empty-state-description">
                            {searchQuery
                                ? 'Try a different search term'
                                : 'Contacts will be automatically created when you receive messages'
                            }
                        </div>
                    </div>
                ) : (
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ borderBottom: '1px solid var(--color-border)' }}>
                                <th style={{ textAlign: 'left', padding: 'var(--space-md)', color: 'var(--color-text-secondary)', fontWeight: 500 }}>
                                    Contact
                                </th>
                                <th style={{ textAlign: 'left', padding: 'var(--space-md)', color: 'var(--color-text-secondary)', fontWeight: 500 }}>
                                    Phone
                                </th>
                                <th style={{ textAlign: 'left', padding: 'var(--space-md)', color: 'var(--color-text-secondary)', fontWeight: 500 }}>
                                    Platforms
                                </th>
                                <th style={{ textAlign: 'left', padding: 'var(--space-md)', color: 'var(--color-text-secondary)', fontWeight: 500 }}>
                                    Added
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {filteredContacts.map((contact) => (
                                <tr
                                    key={contact.id}
                                    style={{ borderBottom: '1px solid var(--color-divider)', cursor: 'pointer' }}
                                >
                                    <td style={{ padding: 'var(--space-md)' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-md)' }}>
                                            <div className="conversation-avatar" style={{ width: '40px', height: '40px', fontSize: 'var(--text-sm)' }}>
                                                {getInitials(contact.name)}
                                            </div>
                                            <div>
                                                <div style={{ fontWeight: 500 }}>{contact.name}</div>
                                                {contact.email && (
                                                    <div style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)' }}>
                                                        {contact.email}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    </td>
                                    <td style={{ padding: 'var(--space-md)', color: 'var(--color-text-secondary)' }}>
                                        {contact.phone || '-'}
                                    </td>
                                    <td style={{ padding: 'var(--space-md)' }}>
                                        <div style={{ display: 'flex', gap: 'var(--space-xs)' }}>
                                            {contact.whatsapp_id && (
                                                <span className="platform-badge whatsapp">WA</span>
                                            )}
                                            {contact.instagram_id && (
                                                <span className="platform-badge instagram">IG</span>
                                            )}
                                        </div>
                                    </td>
                                    <td style={{ padding: 'var(--space-md)', color: 'var(--color-text-secondary)' }}>
                                        {new Date(contact.created_at).toLocaleDateString('id-ID')}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
}
