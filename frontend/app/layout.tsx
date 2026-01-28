import './globals.css';
import type { Metadata } from 'next';
import Sidebar from '@/components/layout/Sidebar';

export const metadata: Metadata = {
    title: 'Omnichannel Dashboard',
    description: 'WhatsApp & Instagram messaging dashboard',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="id">
            <body>
                <div className="app-layout">
                    <Sidebar />
                    <main className="main-content">
                        {children}
                    </main>
                </div>
            </body>
        </html>
    );
}
