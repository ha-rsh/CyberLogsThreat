'use client';

import { useRouter, usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Shield, LogOut } from 'lucide-react';

export function Navbar() {
  const router = useRouter();
  const pathname = usePathname();
  const [username, setUsername] = useState<string | null>(null);
  const [isClient, setIsClient] = useState(false);

  // Load username after component mounts (client-side only)
  useEffect(() => {
    setIsClient(true);
    setUsername(localStorage.getItem('username'));
  }, []);

  // Also update username when pathname changes (after login/register)
  useEffect(() => {
    if (isClient) {
      setUsername(localStorage.getItem('username'));
    }
  }, [pathname, isClient]);

  if (pathname === '/login' || pathname === '/register') {
    return null;
  }

  const handleLogout = () => {
    localStorage.clear();
    sessionStorage.clear();
    window.location.href = '/login';
  };

  return (
    <nav className="border-b bg-white dark:bg-slate-950 sticky top-0 z-50 shadow-sm">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
            <Shield className="h-8 w-8 text-blue-600" />
            <div>
              <h1 className="text-xl font-bold">Cyber Threat Detection System</h1>
            </div>
          </Link>

          {isClient && username && (
            <div className="flex items-center gap-4">
              <span className="text-sm text-muted-foreground hidden md:block">
                Welcome, <strong>{username}</strong>
              </span>
              <Button onClick={handleLogout} variant="outline" size="sm">
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}