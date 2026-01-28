'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navItems = [
    { href: '/dashboard', label: 'Dashboard', icon: '📊' },
    { href: '/inbox', label: 'Inbox', icon: '💬' },
    { href: '/contacts', label: 'Contacts', icon: '👥' },
    { href: '/settings', label: 'Settings', icon: '⚙️' },
];

export default function Sidebar() {
    const pathname = usePathname();

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <div className="sidebar-logo">O</div>
                <span className="sidebar-title">Omnichannel</span>
            </div>

            <nav className="sidebar-nav">
                {navItems.map((item) => (
                    <Link
                        key={item.href}
                        href={item.href}
                        className={`nav-item ${pathname === item.href ? 'active' : ''}`}
                    >
                        <span className="nav-item-icon">{item.icon}</span>
                        <span>{item.label}</span>
                    </Link>
                ))}
            </nav>

            <div style={{ padding: 'var(--space-md)', borderTop: '1px solid var(--color-divider)' }}>
                <div style={{ fontSize: 'var(--text-xs)', color: 'var(--color-text-muted)' }}>
                    Connected Platforms
                </div>
                <div style={{ display: 'flex', gap: 'var(--space-sm)', marginTop: 'var(--space-sm)' }}>
                    <span className="platform-badge whatsapp">WA</span>
                    <span className="platform-badge instagram">IG</span>
                </div>
            </div>
        </aside>
    );
}
