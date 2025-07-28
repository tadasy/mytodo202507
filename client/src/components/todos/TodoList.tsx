import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import type { Todo } from '@/types/api';
import { TodoItem } from './TodoItem';
import { TodoForm } from './TodoForm';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Plus, CheckCircle, Circle } from 'lucide-react';

export function TodoList() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCompleted, setShowCompleted] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    loadTodos();
  }, []);

  const loadTodos = async () => {
    try {
      setIsLoading(true);
      const todosList = await apiClient.getTodos();
      setTodos(todosList || []);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Todoの取得に失敗しました';
      setTodos([]); // エラー時は空配列を設定
      toast({
        title: "エラー",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleTodoCreated = (newTodo: Todo) => {
    setTodos(prev => [newTodo, ...(prev || [])]);
    setShowForm(false);
    toast({
      title: "Todo作成完了",
      description: "新しいTodoが作成されました",
    });
  };

  const handleTodoUpdated = (updatedTodo: Todo) => {
    setTodos(prev => (prev || []).map(todo => 
      todo.id === updatedTodo.id ? updatedTodo : todo
    ));
    toast({
      title: "Todo更新完了",
      description: "Todoが更新されました",
    });
  };

  const handleTodoDeleted = (todoId: string) => {
    setTodos(prev => (prev || []).filter(todo => todo.id !== todoId));
    toast({
      title: "Todo削除完了",
      description: "Todoが削除されました",
    });
  };

  const handleTodoCompleted = async (todoId: string, completed: boolean) => {
    try {
      const updatedTodo = await apiClient.markTodoComplete(todoId, {
        completed,
      });
      setTodos(prev => (prev || []).map(todo => 
        todo.id === todoId ? updatedTodo : todo
      ));
      toast({
        title: completed ? "Todo完了" : "Todo未完了に戻しました",
        description: completed ? "お疲れ様でした！" : "引き続き頑張りましょう",
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Todo状態の更新に失敗しました';
      toast({
        title: "エラー",
        description: errorMessage,
        variant: "destructive",
      });
    }
  };

  const filteredTodos = (todos || []).filter(todo => 
    showCompleted ? todo.completed : !todo.completed
  );

  const completedCount = (todos || []).filter(todo => todo.completed).length;
  const activeCount = (todos || []).filter(todo => !todo.completed).length;

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Todoを読み込み中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Todo管理</h1>
        <div className="flex items-center gap-4 mb-4">
          <Badge variant="secondary" className="flex items-center gap-2">
            <Circle className="h-4 w-4" />
            未完了: {activeCount}
          </Badge>
          <Badge variant="outline" className="flex items-center gap-2">
            <CheckCircle className="h-4 w-4" />
            完了済み: {completedCount}
          </Badge>
        </div>
        
        <div className="flex items-center gap-4">
          <Button
            onClick={() => setShowForm(!showForm)}
            className="flex items-center gap-2"
          >
            <Plus className="h-4 w-4" />
            新しいTodo
          </Button>
          
          <div className="flex gap-2">
            <Button
              variant={!showCompleted ? "default" : "outline"}
              onClick={() => setShowCompleted(false)}
              size="sm"
            >
              未完了
            </Button>
            <Button
              variant={showCompleted ? "default" : "outline"}
              onClick={() => setShowCompleted(true)}
              size="sm"
            >
              完了済み
            </Button>
          </div>
        </div>
      </div>

      {showForm && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>新しいTodo</CardTitle>
            <CardDescription>
              新しいTodoを作成します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TodoForm
              onSubmit={handleTodoCreated}
              onCancel={() => setShowForm(false)}
            />
          </CardContent>
        </Card>
      )}

      <div className="space-y-4">
        {filteredTodos.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <div className="text-gray-500">
                {showCompleted ? (
                  <div>
                    <CheckCircle className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                    <p>完了済みのTodoはありません</p>
                  </div>
                ) : (
                  <div>
                    <Circle className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                    <p>未完了のTodoはありません</p>
                    <p className="text-sm mt-2">新しいTodoを作成してください</p>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        ) : (
          filteredTodos.map(todo => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onUpdate={handleTodoUpdated}
              onDelete={handleTodoDeleted}
              onToggleComplete={handleTodoCompleted}
            />
          ))
        )}
      </div>
    </div>
  );
}
