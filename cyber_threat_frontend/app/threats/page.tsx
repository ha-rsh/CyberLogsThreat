'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useGetThreatsQuery, useAnalyzeThreatsMutation, useLazySearchThreatsQuery } from '@/lib/features/api/apiSlice';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ArrowLeft, Search, Play, Loader2, AlertTriangle } from 'lucide-react';
import { format } from 'date-fns';
import { toast } from 'sonner';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

export default function ThreatsPage() {
  const { data, isLoading, error, refetch } = useGetThreatsQuery();
  const [analyzeThreats, { isLoading: isAnalyzing }] = useAnalyzeThreatsMutation();
  const [searchThreats, { data: searchData, isLoading: isSearching }] = useLazySearchThreatsQuery();
  
  const [filters, setFilters] = useState({
    type: '',
    user: '',
  });

  const handleAnalyze = async () => {
    try {
      const result = await analyzeThreats().unwrap();
      toast.success('Analysis Complete', {
        description: `Detected ${result.data.threatsDetected} threats in ${result.data.duration}`,
      });
      refetch();
    } catch (err: any) {
      toast.error('Analysis Failed', {
        description: err.data?.error?.message || 'An error occurred',
      });
    }
  };

  const handleSearch = () => {
    const searchFilters = {
      type: filters.type === 'all' ? '' : filters.type,
      user: filters.user,
    };
    searchThreats(searchFilters);
  };

  const handleClear = () => {
    setFilters({ type: '', user: '' });
  };

  const threats = searchData?.data || data?.data || [];
  const loading = isLoading || isSearching;

  const getSeverityBadge = (severity: string) => {
    const variants: Record<string, { variant: 'default' | 'secondary' | 'destructive' | 'outline', className: string }> = {
      Critical: { variant: 'destructive', className: 'bg-red-600' },
      High: { variant: 'destructive', className: 'bg-orange-600' },
      Medium: { variant: 'outline', className: 'border-yellow-600 text-yellow-600' },
      Low: { variant: 'secondary', className: '' },
    };
    const config = variants[severity] || variants.Low;
    return (
      <Badge variant={config.variant} className={config.className}>
        {severity}
      </Badge>
    );
  };

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-950">
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-4">
            <Link href="/">
              <Button variant="ghost" size="icon">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
            <div>
              <h1 className="text-3xl font-bold">Detected Threats</h1>
              <p className="text-muted-foreground">View and analyze security threats</p>
            </div>
          </div>
          <Button onClick={handleAnalyze} disabled={isAnalyzing}>
            {isAnalyzing ? (
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            ) : (
              <Play className="h-4 w-4 mr-2" />
            )}
            {isAnalyzing ? 'Analyzing...' : 'Run Analysis'}
          </Button>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Search Filters</CardTitle>
            <CardDescription>Filter threats by type or user</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <Input
                  id="user"
                  placeholder="Search with userId"
                  value={filters.user}
                  onChange={(e) => setFilters({ ...filters, user: e.target.value })}
                />
              </div>
              <div>
                <Select 
                  value={filters.type} 
                  onValueChange={(value) => setFilters({ ...filters, type: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select threat type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="Credential Stuffing">Credential Stuffing</SelectItem>
                    <SelectItem value="Privilege Escalation">Privilege Escalation</SelectItem>
                    <SelectItem value="Account Takeover">Account Takeover</SelectItem>
                    <SelectItem value="Data Exfiltration">Data Exfiltration</SelectItem>
                    <SelectItem value="Insider Threat">Insider Threat</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex items-end gap-2">
                <Button onClick={handleSearch} className="flex-1">
                  <Search className="h-4 w-4 mr-2" />
                  Search
                </Button>
                <Button onClick={handleClear} variant="outline">
                  Clear
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-red-600" />
              Threats ({threats.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertDescription>
                  Failed to load threats. Please check your connection.
                </AlertDescription>
              </Alert>
            )}

            {loading ? (
              <div className="flex justify-center items-center py-8">
                <Loader2 className="h-8 w-8 animate-spin" />
              </div>
            ) : threats.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No threats detected. Run analysis to scan for threats.
              </div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Timestamp</TableHead>
                      <TableHead>User ID</TableHead>
                      <TableHead>IP Address</TableHead>
                      <TableHead>Threat Type</TableHead>
                      <TableHead>Severity</TableHead>
                      <TableHead>Action</TableHead>
                      <TableHead>Details</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {threats.map((threat: any, idx: number) => (
                      <TableRow key={threat.id || idx}>
                        <TableCell className="font-mono text-xs">
                          {format(new Date(threat.timestamp), 'yyyy-MM-dd HH:mm:ss')}
                        </TableCell>
                        <TableCell className="font-mono">{threat.userId}</TableCell>
                        <TableCell className="font-mono text-xs">{threat.ipAddress}</TableCell>
                        <TableCell>
                          <Badge variant="outline">{threat.threatType}</Badge>
                        </TableCell>
                        <TableCell>{getSeverityBadge(threat.severity)}</TableCell>
                        <TableCell>
                          <Badge variant="secondary">{threat.action}</Badge>
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {threat.fileName && <div>File: {threat.fileName}</div>}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}