import axios from 'axios';
import { toast } from '@/hooks/use-toast';
import { Response } from '@/domain/common';
import { useAuthStore } from '@/store/authStore';



const forbiddenCode = 3999;

// 创建 axios 实例
let url = import.meta.env.VITE_API_BASE_URL || '/api';
url = "http://localhost:1314"
const axiosInstance = axios.create({
  baseURL: url,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
axiosInstance.interceptors.request.use(
  (config) => {
    // 可以在这里添加 token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器'
axiosInstance.interceptors.response.use(
  (response) => {
    const data = response.data as Response<any>;

    // 检查 code 为 3999，跳转到登录页
    if (data.code === forbiddenCode) {
      toast({
        title: '未授权',
        description: '请重新登录',
        variant: 'destructive',
      });
      useAuthStore.getState().logout();
      window.location.href = '/login';
      return Promise.reject(new Error('Unauthorized'));
    }

    if (data.code !== 0) {
      toast({
        title: '错误',
        description: data.message,
        variant: 'destructive',
      });
      return Promise.reject(new Error(data.message));
    }

    return response;
  },
  (error) => {
    // 统一错误处理
    const message = error.response?.data?.message || error.message || '请求失败';

    toast({
      title: '错误',
      description: message,
      variant: 'destructive',
    });

    // 401 未授权或 code 为 3999，跳转到登录页
    if (error.response?.status === 401 || error.response?.data?.code === forbiddenCode) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }


    return Promise.reject(error);
  }
);

export default axiosInstance;
