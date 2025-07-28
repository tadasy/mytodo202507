import { useRouteError, isRouteErrorResponse } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertTriangle, RefreshCw } from 'lucide-react';

export function ErrorBoundary() {
  const error = useRouteError();

  const handleReload = () => {
    window.location.reload();
  };

  const handleGoHome = () => {
    window.location.href = '/';
  };

  if (isRouteErrorResponse(error)) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 p-4">
        <Card className="w-full max-w-md">
          <CardHeader>
            <div className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-red-600" />
              <CardTitle>エラーが発生しました</CardTitle>
            </div>
            <CardDescription>
              ページの読み込み中にエラーが発生しました
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Alert variant="destructive">
              <AlertTitle>HTTP {error.status}</AlertTitle>
              <AlertDescription>
                {error.statusText || error.data?.message || 'リクエストの処理中にエラーが発生しました'}
              </AlertDescription>
            </Alert>
            
            <div className="flex gap-2">
              <Button onClick={handleReload} variant="outline" className="flex-1">
                <RefreshCw className="h-4 w-4 mr-2" />
                再読み込み
              </Button>
              <Button onClick={handleGoHome} className="flex-1">
                ホームへ戻る
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <div className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-red-600" />
            <CardTitle>予期しないエラー</CardTitle>
          </div>
          <CardDescription>
            アプリケーションで予期しないエラーが発生しました
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Alert variant="destructive">
            <AlertTitle>エラー詳細</AlertTitle>
            <AlertDescription>
              {error instanceof Error ? error.message : 'Unknown error occurred'}
            </AlertDescription>
          </Alert>
          
          {/* 開発モードでのみスタックトレースを表示 */}
          {error instanceof Error && error.stack && (
            <details className="mt-4">
              <summary className="cursor-pointer text-sm font-medium mb-2">
                開発者向け詳細情報
              </summary>
              <pre className="text-xs bg-gray-100 p-2 rounded overflow-auto max-h-40">
                {error.stack}
              </pre>
            </details>
          )}
          
          <div className="flex gap-2">
            <Button onClick={handleReload} variant="outline" className="flex-1">
              <RefreshCw className="h-4 w-4 mr-2" />
              再読み込み
            </Button>
            <Button onClick={handleGoHome} className="flex-1">
              ホームへ戻る
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
