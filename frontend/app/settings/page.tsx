export default function SettingsPage() {
    return (
        <div className="page-container">
            <h1 style={{ fontSize: 'var(--text-2xl)', fontWeight: 700, marginBottom: 'var(--space-lg)' }}>
                Settings
            </h1>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-lg)', maxWidth: '600px' }}>
                {/* WhatsApp Settings */}
                <div className="card">
                    <div className="card-header">
                        <h2 className="card-title">WhatsApp Business</h2>
                        <span className="platform-badge whatsapp">Connected</span>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-md)' }}>
                        <div>
                            <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                                Phone Number ID
                            </label>
                            <input
                                type="text"
                                className="search-input"
                                placeholder="Enter Phone Number ID"
                                style={{ paddingLeft: 'var(--space-md)' }}
                            />
                        </div>
                        <div>
                            <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                                Business Account ID
                            </label>
                            <input
                                type="text"
                                className="search-input"
                                placeholder="Enter Business Account ID"
                                style={{ paddingLeft: 'var(--space-md)' }}
                            />
                        </div>
                    </div>
                </div>

                {/* Instagram Settings */}
                <div className="card">
                    <div className="card-header">
                        <h2 className="card-title">Instagram</h2>
                        <span className="platform-badge instagram">Connected</span>
                    </div>
                    <div>
                        <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                            Instagram Account ID
                        </label>
                        <input
                            type="text"
                            className="search-input"
                            placeholder="Enter Instagram Account ID"
                            style={{ paddingLeft: 'var(--space-md)' }}
                        />
                    </div>
                </div>

                {/* Webhook Settings */}
                <div className="card">
                    <div className="card-header">
                        <h2 className="card-title">Webhooks</h2>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-md)' }}>
                        <div>
                            <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                                WhatsApp Webhook URL
                            </label>
                            <div style={{
                                padding: 'var(--space-sm) var(--space-md)',
                                backgroundColor: 'var(--color-bg-tertiary)',
                                borderRadius: 'var(--radius-md)',
                                fontFamily: 'var(--font-mono)',
                                fontSize: 'var(--text-sm)',
                                color: 'var(--color-text-secondary)'
                            }}>
                                https://omni.otomasi.click/webhooks/whatsapp
                            </div>
                        </div>
                        <div>
                            <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                                Instagram Webhook URL
                            </label>
                            <div style={{
                                padding: 'var(--space-sm) var(--space-md)',
                                backgroundColor: 'var(--color-bg-tertiary)',
                                borderRadius: 'var(--radius-md)',
                                fontFamily: 'var(--font-mono)',
                                fontSize: 'var(--text-sm)',
                                color: 'var(--color-text-secondary)'
                            }}>
                                https://omni.otomasi.click/webhooks/instagram
                            </div>
                        </div>
                        <div>
                            <label style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', display: 'block', marginBottom: 'var(--space-xs)' }}>
                                Verify Token
                            </label>
                            <div style={{
                                padding: 'var(--space-sm) var(--space-md)',
                                backgroundColor: 'var(--color-bg-tertiary)',
                                borderRadius: 'var(--radius-md)',
                                fontFamily: 'var(--font-mono)',
                                fontSize: 'var(--text-sm)',
                                color: 'var(--color-text-secondary)'
                            }}>
                                omnichannel_verify_token
                            </div>
                        </div>
                    </div>
                </div>

                <button className="btn btn-primary" style={{ alignSelf: 'flex-start' }}>
                    Save Settings
                </button>
            </div>
        </div>
    );
}
