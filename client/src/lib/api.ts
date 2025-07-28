import {
  User,
  Todo,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  CreateTodoRequest,
  UpdateTodoRequest,
  MarkTodoCompleteRequest,
} from '../types/api';

const API_BASE_URL = 'http://localhost:8080/api';

class ApiClient {
  private token: string | null = null;

  constructor() {
    // トークンをlocalStorageから取得
    this.token = localStorage.getItem('auth_token');
  }

  private getAuthHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    return headers;
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // 認証API
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    const result = await this.handleResponse<AuthResponse>(response);
    this.setTokenInternal(result.token);
    return result;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    const result = await this.handleResponse<AuthResponse>(response);
    this.setTokenInternal(result.token);
    return result;
  }

  logout(): void {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  private setTokenInternal(token: string): void {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  public setToken(token: string): void {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  public getProfile(): Promise<User> {
    // プロフィール取得はユーザーサービスから実装が必要
    // 現在は実装されていないため、ダミーのユーザー情報を返す
    return Promise.resolve({
      id: '1',
      email: 'user@example.com',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    });
  }

  // Todo API
  async getTodos(): Promise<Todo[]> {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse<Todo[]>(response);
  }

  async getCompletedTodos(): Promise<Todo[]> {
    const response = await fetch(`${API_BASE_URL}/todos?completed=true`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse<Todo[]>(response);
  }

  async getTodo(id: string): Promise<Todo> {
    const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse<Todo>(response);
  }

  async createTodo(data: CreateTodoRequest): Promise<Todo> {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    return this.handleResponse<Todo>(response);
  }

  async updateTodo(id: string, data: UpdateTodoRequest): Promise<Todo> {
    const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    return this.handleResponse<Todo>(response);
  }

  async markTodoComplete(id: string, data: MarkTodoCompleteRequest): Promise<Todo> {
    const response = await fetch(`${API_BASE_URL}/todos/${id}/complete`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    return this.handleResponse<Todo>(response);
  }

  async deleteTodo(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });

    await this.handleResponse<{ message: string }>(response);
  }

  // トークンの存在チェック
  isAuthenticated(): boolean {
    return !!this.token;
  }
}

// シングルトンインスタンス
export const apiClient = new ApiClient();
