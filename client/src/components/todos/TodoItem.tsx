import { useState } from 'react';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import type { Todo } from '@/types/api';
import { TodoForm } from './TodoForm';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { Edit2, Trash2, Calendar } from 'lucide-react';

interface TodoItemProps {
  todo: Todo;
  onUpdate: (todo: Todo) => void;
  onDelete: (todoId: string) => void;
  onToggleComplete: (todoId: string, completed: boolean) => void;
}

export function TodoItem({ todo, onUpdate, onDelete, onToggleComplete }: TodoItemProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const { toast } = useToast();

  const handleDelete = async () => {
    if (!window.confirm('このTodoを削除しますか？')) {
      return;
    }

    setIsDeleting(true);
    try {
      await apiClient.deleteTodo(todo.id);
      onDelete(todo.id);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Todo削除に失敗しました';
      toast({
        title: "エラー",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const handleToggleComplete = () => {
    onToggleComplete(todo.id, !todo.completed);
  };

  const handleUpdate = (updatedTodo: Todo) => {
    onUpdate(updatedTodo);
    setIsEditing(false);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isEditing) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="mb-4">
            <h3 className="text-lg font-semibold">Todoを編集</h3>
          </div>
          <TodoForm
            todo={todo}
            onSubmit={handleUpdate}
            onCancel={() => setIsEditing(false)}
          />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={`transition-all ${todo.completed ? 'opacity-75 bg-gray-50' : ''}`}>
      <CardContent className="pt-6">
        <div className="flex items-start gap-4">
          <div className="flex items-center pt-1">
            <Checkbox
              checked={todo.completed}
              onCheckedChange={handleToggleComplete}
              disabled={isDeleting}
            />
          </div>
          
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <h3 className={`text-lg font-semibold ${
                  todo.completed ? 'line-through text-gray-600' : 'text-gray-900'
                }`}>
                  {todo.title}
                </h3>
                
                {todo.description && (
                  <p className={`mt-2 text-sm ${
                    todo.completed ? 'line-through text-gray-500' : 'text-gray-700'
                  }`}>
                    {todo.description}
                  </p>
                )}
                
                <div className="flex items-center gap-4 mt-4 text-xs text-gray-500">
                  <div className="flex items-center gap-1">
                    <Calendar className="h-3 w-3" />
                    作成: {formatDate(todo.created_at)}
                  </div>
                  {todo.updated_at !== todo.created_at && (
                    <div className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      更新: {formatDate(todo.updated_at)}
                    </div>
                  )}
                </div>
              </div>
              
              <div className="flex flex-col items-end gap-2">
                <div className="flex items-center gap-2">
                  <Badge variant={todo.completed ? "outline" : "secondary"}>
                    {todo.completed ? "完了" : "未完了"}
                  </Badge>
                </div>
                
                <div className="flex items-center gap-1">
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => setIsEditing(true)}
                    disabled={isDeleting}
                    className="h-8 w-8 p-0"
                  >
                    <Edit2 className="h-4 w-4" />
                  </Button>
                  
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={handleDelete}
                    disabled={isDeleting}
                    className="h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
