'use client';

import Link from 'next/link';
import { useState } from 'react';
import { useGetLogsQuery, useGetThreatsQuery } from '@/lib/features/api/apiSlice';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Shield, FileText, AlertTriangle, Upload, Loader2, Activity } from 'lucide-react';
import { BarChart, Bar, LineChart, Line, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useRouter } from 'next/navigation';

const COLORS = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899'];

export default function Home() {
  const router = useRouter();
  const { data: logsData, isLoading: logsLoading } = useGetLogsQuery();
  const { data: threatsData, isLoading: threatsLoading } = useGetThreatsQuery();

  const logs = logsData?.data || [];
  const threats = threatsData?.data || [];
  const loading = logsLoading || threatsLoading;

  const totalLogs = logs.length;
  const totalThreats = threats.length;
  const uniqueUsers = new Set(logs.map((log: any) => log.userId)).size;
  const threatPercentage = totalLogs > 0 ? ((totalThreats / totalLogs) * 100).toFixed(2) : '0';

  // Action distribution
  const actionCounts = logs.reduce((acc: any, log: any) => {
    acc[log.action] = (acc[log.action] || 0) + 1;
    return acc;
  }, {});
  const actionData = Object.entries(actionCounts).map(([name, value]) => ({ name, value }));

  const threatTypeCounts = threats.reduce((acc: any, threat: any) => {
    acc[threat.threatType] = (acc[threat.threatType] || 0) + 1;
    return acc;
  }, {});
  const threatTypeData = Object.entries(threatTypeCounts).map(([name, value]) => ({ name, value }));

  const severityCounts = threats.reduce((acc: any, threat: any) => {
    acc[threat.severity] = (acc[threat.severity] || 0) + 1;
    return acc;
  }, {});
  const severityData = Object.entries(severityCounts).map(([name, value]) => ({ name, value }));

  const timelineData = logs.reduce((acc: any, log: any) => {
    const hour = new Date(log.timestamp).getHours();
    const key = `${hour}:00`;
    acc[key] = (acc[key] || 0) + 1;
    return acc;
  }, {});
  const timeline = Object.entries(timelineData)
    .map(([time, count]) => ({ time, count }))
    .sort((a, b) => parseInt(a.time) - parseInt(b.time));

  const handleLogout = () => {
    localStorage.clear();
    sessionStorage.clear();
    window.location.href = '/login';
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900">
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-8">
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <Upload className="h-8 w-8 mb-2 text-blue-600" />
              <CardTitle>Upload Logs</CardTitle>
              <CardDescription>Import CSV log files</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/logs/upload">
                <Button className="w-full">Upload</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <FileText className="h-8 w-8 mb-2 text-green-600" />
              <CardTitle>View Logs</CardTitle>
              <CardDescription>Browse all system logs</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/logs">
                <Button className="w-full" variant="outline">View Logs</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <AlertTriangle className="h-8 w-8 mb-2 text-red-600" />
              <CardTitle>Threats</CardTitle>
              <CardDescription>View detected threats</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/threats">
                <Button className="w-full" variant="outline">View Threats</Button>
              </Link>
            </CardContent>
          </Card>
        </div>

        {loading ? (
          <div className="flex justify-center items-center py-16">
            <Loader2 className="h-12 w-12 animate-spin" />
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Total Logs</CardTitle>
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{totalLogs}</div>
                  <p className="text-xs text-muted-foreground">System events recorded</p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Total Threats</CardTitle>
                  <AlertTriangle className="h-4 w-4 text-red-600" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-red-600">{totalThreats}</div>
                  <p className="text-xs text-muted-foreground">Security threats detected</p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Unique Users</CardTitle>
                  <Activity className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{uniqueUsers}</div>
                  <p className="text-xs text-muted-foreground">Active users monitored</p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Threat Rate</CardTitle>
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{threatPercentage}%</div>
                  <p className="text-xs text-muted-foreground">Threats per total logs</p>
                </CardContent>
              </Card>
            </div>

            {/* Charts Section */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
              <Card>
                <CardHeader>
                  <CardTitle>Action Distribution</CardTitle>
                  <CardDescription>Breakdown of log actions</CardDescription>
                </CardHeader>
                <CardContent>
                  {actionData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <PieChart>
                        <Pie
                          data={actionData}
                          cx="50%"
                          cy="50%"
                          labelLine={false}
                          label={({ name, percent }: any) => `${name}: ${(percent * 100).toFixed(0)}%`}
                          outerRadius={80}
                          fill="#8884d8"
                          dataKey="value"
                        >
                          {actionData.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                          ))}
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="text-center py-8 text-muted-foreground">No data available</div>
                  )}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Threat Types</CardTitle>
                  <CardDescription>Distribution of detected threats</CardDescription>
                </CardHeader>
                <CardContent>
                  {threatTypeData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={threatTypeData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} fontSize={12} />
                        <YAxis />
                        <Tooltip />
                        <Bar dataKey="value" fill="#ef4444" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="text-center py-8 text-muted-foreground">No threats detected</div>
                  )}
                </CardContent>
              </Card>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Activity Timeline</CardTitle>
                  <CardDescription>Logs per hour</CardDescription>
                </CardHeader>
                <CardContent>
                  {timeline.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <LineChart data={timeline}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="time" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Line type="monotone" dataKey="count" stroke="#3b82f6" strokeWidth={2} />
                      </LineChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="text-center py-8 text-muted-foreground">No data available</div>
                  )}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Threat Severity</CardTitle>
                  <CardDescription>Severity level breakdown</CardDescription>
                </CardHeader>
                <CardContent>
                  {severityData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={severityData} layout="vertical">
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis type="number" />
                        <YAxis dataKey="name" type="category" width={80} />
                        <Tooltip />
                        <Bar dataKey="value" fill="#f59e0b" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="text-center py-8 text-muted-foreground">No threats detected</div>
                  )}
                </CardContent>
              </Card>
            </div>
          </>
        )}
      </div>
    </div>
  );
}