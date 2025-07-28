import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { LoginForm } from './components/auth/LoginForm';
import { RegisterForm } from './components/auth/RegisterForm';
import { TodoList } from './components/todos/TodoList';
import { DashboardLayout } from './components/layout/DashboardLayout';
import { Toaster } from './components/ui/toaster';
import { ErrorBoundary } from './components/ErrorBoundary';
import './App.css';

// React Router v7のルーター設定
const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginForm />,
    errorElement: <ErrorBoundary />,
  },
  {
    path: '/register', 
    element: <RegisterForm />,
    errorElement: <ErrorBoundary />,
  },
  {
    path: '/todos',
    element: (
      <ProtectedRoute>
        <DashboardLayout>
          <TodoList />
        </DashboardLayout>
      </ProtectedRoute>
    ),
    errorElement: <ErrorBoundary />,
  },
  {
    path: '/',
    element: <Navigate to="/todos" replace />,
    errorElement: <ErrorBoundary />,
  },
  {
    path: '*',
    element: <Navigate to="/todos" replace />,
    errorElement: <ErrorBoundary />,
  },
]);

function App() {
  return (
    <AuthProvider>
      <div className="App">
        <RouterProvider router={router} />
        <Toaster />
      </div>
    </AuthProvider>
  );
}

export default App;
