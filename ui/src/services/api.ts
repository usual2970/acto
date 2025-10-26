import { Response } from '@/domain/common';
import { PointType } from '@/domain/point_type';
import axios from '@/lib/axios';



type LoginResponse = {
  token: string;
}

// 积分类型相关接口
export const pointsTypeApi = {
  // 获取积分类型列表
  getList: async () => {
    const resp = await axios.get<Response<PointType[]>>('/admin/v1/point-types');
    if (resp.data.code !== 0) {
      throw new Error(resp.data.message);
    }
    return resp.data.data;
  },

  // 获取积分类型详情
  getDetail: async (id: string) => {
    const resp = await axios.get<Response<PointType>>(`/admin/v1/point-types/${id}`);
    if (resp.data.code !== 0) {
      throw new Error(resp.data.message);
    }
    return resp.data.data;
  },

  // 创建积分类型
  create: async (data: any) => {
    const resp = await axios.post<Response<PointType>>('/admin/v1/point-types', data);
    if (resp.data.code !== 0) {
      throw new Error(resp.data.message);
    }
    return resp.data.data;
  },

  // 更新积分类型
  update: async (id: string, data: any) => {
    const resp = await axios.put<Response<PointType>>(`/admin/v1/point-types/${id}`, data);
    if (resp.data.code !== 0) {
      throw new Error(resp.data.message);
    }
    return resp.data.data;
  },

  // 删除积分类型
  delete: (id: string) => axios.delete(`/admin/v1/point-types/${id}`),
};

// 用户积分相关接口
export const userPointsApi = {
  // 获取用户积分列表
  getList: (pointsTypeId: string) => axios.get(`/points-types/${pointsTypeId}/user-points`),

  // 调整用户积分
  adjust: (pointsTypeId: string, data: any) => axios.post(`/points-types/${pointsTypeId}/user-points/adjust`, data),
};

// 排行榜相关接口
export const leaderboardApi = {
  // 获取排行榜
  getList: (pointsTypeId: string) => axios.get(`/points-types/${pointsTypeId}/leaderboard`),
};

// 奖励相关接口
export const rewardApi = {
  // 获取奖励列表
  getList: (pointsTypeId: string) => axios.get(`/points-types/${pointsTypeId}/rewards`),

  // 获取奖励记录
  getRecords: (pointsTypeId: string) => axios.get(`/points-types/${pointsTypeId}/reward-records`),

  // 创建奖励
  create: (pointsTypeId: string, data: any) => axios.post(`/points-types/${pointsTypeId}/rewards`, data),

  // 更新奖励
  update: (pointsTypeId: string, rewardId: string, data: any) =>
    axios.put(`/points-types/${pointsTypeId}/rewards/${rewardId}`, data),

  // 删除奖励
  delete: (pointsTypeId: string, rewardId: string) =>
    axios.delete(`/points-types/${pointsTypeId}/rewards/${rewardId}`),
};

// 登录相关接口
export const authApi = {
  // 登录
  login: async (data: { username: string; password: string }) => {
    const resp = await axios.post<Response<LoginResponse>>('/admin/v1/login', data);
    if (resp.data.code !== 0) {
      throw new Error(resp.data.message);
    }

    return resp.data.data;
  },

  // 登出
  logout: () => axios.post('/auth/logout'),
};
