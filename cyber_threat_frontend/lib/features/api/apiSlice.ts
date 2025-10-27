import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const getToken = () => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
};

export interface Log {
  id?: string;
  timestamp: string;
  userId: string;
  ipAddress: string;
  action: string;
  fileName?: string;
  databaseQuery?: string;
}

export interface Threat {
  id?: string;
  timestamp: string;
  userId: string;
  ipAddress: string;
  action: string;
  fileName?: string;
  threatType: string;
  severity: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: number;
    message: string;
    details?: string;
  };
  meta?: {
    count?: number;
    total?: number;
    timestamp: number;
  };
}

export interface AnalysisResult {
  threatsDetected: number;
  duration: string;
}

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
    prepareHeaders: (headers) => {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['Logs', 'Threats'],
  endpoints: (builder) => ({
    getLogs: builder.query<APIResponse<Log[]>, void>({
      query: () => '/api/logs',
      providesTags: ['Logs'],
    }),
    
    getLogById: builder.query<APIResponse<Log>, string>({
      query: (id) => `/api/logs/${id}`,
      providesTags: ['Logs'],
    }),
    
    createLog: builder.mutation<APIResponse<Log>, Partial<Log>>({
      query: (log) => ({
        url: '/api/logs',
        method: 'POST',
        body: log,
      }),
      invalidatesTags: ['Logs'],
    }),
    
    searchLogs: builder.query<APIResponse<Log[]>, { userId?: string; action?: string; startTime?: string; endTime?: string }>({
      query: (params) => ({
        url: '/api/logs/search',
        params,
      }),
      providesTags: ['Logs'],
    }),
    
    getThreats: builder.query<APIResponse<Threat[]>, void>({
      query: () => '/api/threats',
      providesTags: ['Threats'],
    }),
    
    getThreatById: builder.query<APIResponse<Threat>, string>({
      query: (id) => `/api/threats/${id}`,
      providesTags: ['Threats'],
    }),
    
    analyzeThreats: builder.mutation<APIResponse<AnalysisResult>, void>({
      query: () => ({
        url: '/api/threats/analyze',
        method: 'POST',
      }),
      invalidatesTags: ['Threats'],
    }),
    
    searchThreats: builder.query<APIResponse<Threat[]>, { type?: string; user?: string }>({
      query: (params) => ({
        url: '/api/threats/search',
        params,
      }),
      providesTags: ['Threats'],
    }),

    login: builder.mutation({
      query: (credentials) => ({
        url: '/api/auth/login',
        method: 'POST',
        body: credentials,
      }),
    }),
    register: builder.mutation({
      query: (userData) => ({
        url: '/api/auth/register',
        method: 'POST',
        body: userData,
      }),
    }),
    refreshToken: builder.mutation({
      query: () => ({
        url: '/api/auth/refresh',
        method: 'POST',
      }),
    }),
  }),
});

export const {
  useGetLogsQuery,
  useGetLogByIdQuery,
  useCreateLogMutation,
  useSearchLogsQuery,
  useLazySearchLogsQuery,
  useGetThreatsQuery,
  useGetThreatByIdQuery,
  useAnalyzeThreatsMutation,
  useSearchThreatsQuery,
  useLazySearchThreatsQuery,
  useLoginMutation,
  useRegisterMutation,
  useRefreshTokenMutation,
} = apiSlice;