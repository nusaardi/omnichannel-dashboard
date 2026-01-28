const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'https://omni.otomasi.click';

interface FetchOptions extends RequestInit {
    params?: Record<string, string>;
}

class ApiClient {
    private baseUrl: string;

    constructor(baseUrl: string = API_BASE) {
        this.baseUrl = baseUrl;
    }

    private async request<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
        const { params, ...fetchOptions } = options;

        let url = `${this.baseUrl}${endpoint}`;
        if (params) {
            const searchParams = new URLSearchParams(params);
            url += `?${searchParams.toString()}`;
        }

        const response = await fetch(url, {
            ...fetchOptions,
            headers: {
                'Content-Type': 'application/json',
                ...fetchOptions.headers,
            },
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({ error: 'Unknown error' }));
            throw new Error(error.error || `Request failed: ${response.status}`);
        }

        return response.json();
    }

    // Conversations
    async getConversations(limit = 50, offset = 0) {
        return this.request<{ conversations: any[]; total: number }>('/api/conversations', {
            params: { limit: limit.toString(), offset: offset.toString() },
        });
    }

    async getConversation(id: string) {
        return this.request<{ conversation: any; messages: any[] }>(`/api/conversations/${id}`);
    }

    // Messages
    async sendMessage(data: {
        conversation_id?: string;
        platform: string;
        recipient_id: string;
        content: string;
        content_type?: string;
    }) {
        return this.request<any>('/api/messages', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    // Contacts
    async getContacts(limit = 100, offset = 0) {
        return this.request<{ contacts: any[]; total: number }>('/api/contacts', {
            params: { limit: limit.toString(), offset: offset.toString() },
        });
    }

    async createContact(data: { name: string; phone?: string; email?: string }) {
        return this.request<any>('/api/contacts', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    // Health
    async healthCheck() {
        return this.request<{ status: string; timestamp: string }>('/health');
    }
}

export const api = new ApiClient();
export default ApiClient;
