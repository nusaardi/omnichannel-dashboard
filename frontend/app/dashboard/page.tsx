export default function DashboardPage() {
    return (
        <div className="page-container">
            <h1 style={{ fontSize: 'var(--text-2xl)', fontWeight: 700, marginBottom: 'var(--space-lg)' }}>
                Dashboard
            </h1>

            <div className="stats-grid">
                <div className="stat-card">
                    <div className="stat-value">0</div>
                    <div className="stat-label">Messages Today</div>
                    <div className="stat-change positive">
                        <span>↑</span>
                        <span>0% from yesterday</span>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-value">0</div>
                    <div className="stat-label">Active Conversations</div>
                    <div className="stat-change positive">
                        <span>↑</span>
                        <span>0 new</span>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-value">0</div>
                    <div className="stat-label">Total Contacts</div>
                    <div className="stat-change positive">
                        <span>↑</span>
                        <span>0 this week</span>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-value" style={{ color: 'var(--color-whatsapp)' }}>●</div>
                    <div className="stat-label">WhatsApp Status</div>
                    <div className="stat-change" style={{ color: 'var(--color-text-muted)' }}>
                        <span>Connected</span>
                    </div>
                </div>
            </div>

            <div style={{ marginTop: 'var(--space-xl)' }}>
                <div className="card">
                    <div className="card-header">
                        <h2 className="card-title">Recent Activity</h2>
                    </div>
                    <div className="empty-state">
                        <div className="empty-state-icon">📬</div>
                        <div className="empty-state-title">No recent activity</div>
                        <div className="empty-state-description">
                            Messages and conversations will appear here once you start receiving them.
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
