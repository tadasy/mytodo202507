import { useState } from 'react';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import type { Todo } from '@/types/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

interface TodoFormProps {
  todo?: Todo;
  onSubmit: (todo: Todo) => void;
  onCancel: () => void;
}

export function TodoForm({ todo, onSubmit, onCancel }: TodoFormProps) {
  const [title, setTitle] = useState(todo?.title || '');
  const [description, setDescription] = useState(todo?.description || '');
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) {
      toast({
        title: "入力エラー",
        description: "タイトルを入力してください",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);

    try {
      if (todo) {
        // 更新
        const updatedTodo = await apiClient.updateTodo(todo.id, {
          title: title.trim(),
          description: description.trim(),
        });
        onSubmit(updatedTodo);
      } else {
        // 新規作成
        const newTodo = await apiClient.createTodo({
          title: title.trim(),
          description: description.trim(),
        });
        onSubmit(newTodo);
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Todo操作に失敗しました';
      toast({
        title: "エラー",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="title">タイトル *</Label>
        <Input
          id="title"
          type="text"
          placeholder="Todoのタイトルを入力"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          disabled={isLoading}
          maxLength={200}
        />
      </div>
      
      <div className="space-y-2">
        <Label htmlFor="description">説明</Label>
        <Textarea
          id="description"
          placeholder="詳細な説明（任意）"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          disabled={isLoading}
          rows={3}
          maxLength={1000}
        />
      </div>
      
      <div className="flex gap-2 justify-end">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isLoading}
        >
          キャンセル
        </Button>
        <Button
          type="submit"
          disabled={isLoading || !title.trim()}
        >
          {isLoading ? (
            todo ? '更新中...' : '作成中...'
          ) : (
            todo ? '更新' : '作成'
          )}
        </Button>
      </div>
    </form>
  );
}
