import axios, { AxiosInstance } from 'axios';

export interface SDKConfig {
  apiKey: string;
  baseUrl?: string;
}

export class MightyEagle {
  private axios: AxiosInstance;

  constructor(config: SDKConfig) {
    this.axios = axios.create({
      baseURL: config.baseUrl || 'http://localhost:8080/v1',
      headers: {
        'X-API-Key': config.apiKey,
        'Content-Type': 'application/json',
      },
    });
  }

  get persona() {
    return {
      createVerification: async (subjectId: string, provider: string) => {
        const { data } = await this.axios.post('/persona/verifications', {
          subject_id: subjectId,
          provider: provider,
        });
        return data;
      },
      getVerification: async (id: string) => {
        const { data } = await this.axios.get(`/persona/verifications/${id}`);
        return data;
      },
    };
  }

  get reputation() {
    return {
      getScore: async (subjectId: string) => {
        const { data } = await this.axios.get(`/reputation/${subjectId}`);
        return data;
      },
    };
  }

  get consent() {
    return {
      createToken: async (parties: string[], scope: string, expiresAt: string) => {
        const { data } = await this.axios.post('/consent/tokens', {
          parties,
          scope,
          expires_at: expiresAt,
        });
        return data;
      },
      revokeToken: async (id: string, revokedBy: string, reason?: string) => {
        const { data } = await this.axios.post(`/consent/tokens/${id}/revoke`, {
          revoked_by: revokedBy,
          reason,
        });
        return data;
      },
    };
  }

  get webhooks() {
    return {
      listEndpoints: async () => {
        const { data } = await this.axios.get('/webhooks/endpoints');
        return data;
      },
      createEndpoint: async (url: string, events: string[]) => {
        const { data } = await this.axios.post('/webhooks/endpoints', {
          url,
          events,
        });
        return data;
      },
    };
  }

  get billing() {
    return {
      getUsage: async () => {
        const { data } = await this.axios.get('/billing/usage');
        return data;
      },
    };
  }

  get audit() {
    return {
      createExport: async (format: 'csv' | 'json', startDate: string, endDate: string) => {
        const { data } = await this.axios.post('/audit/exports', {
          format,
          start_date: startDate,
          end_date: endDate,
        });
        return data;
      },
      getExport: async (id: string) => {
        const { data } = await this.axios.get(`/audit/exports/${id}`);
        return data;
      },
    };
  }
}
