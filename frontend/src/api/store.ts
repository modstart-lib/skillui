import http from './http';

export interface ApiResponse<T = any> {
  code: number;
  data: T;
  message?: string;
}

export interface SkillItem {
  id: number;
  name: string;
  titleEn?: string;
  titleZh?: string;
  descEn?: string;
  descZh?: string;
  tagsEn?: string[];
  tagsZh?: string[];
  tags?: string[];
  owner?: string;
  version?: string;
  preview?: string;
  [key: string]: any;
}

export interface PaginateParams {
  sort: 'hot' | 'time';
  keywords: string;
  page: number;
  pageSize: number;
}

export interface PaginateResult {
  records: SkillItem[];
  total: number;
}

export const skillUiPaginate = (params: PaginateParams) =>
  http.post<ApiResponse<PaginateResult>>('/skill_ui/paginate', params).then(r => r.data);

export const skillUiDownload = (id: number) =>
  http.post<ApiResponse<{ url: string }>>('/skill_ui/download', { id }).then(r => r.data);

export const skillUiDetail = (id: number) =>
  http.post<ApiResponse<{ record: SkillItem & { _preview?: string } }>>('/skill_ui/get', { id }).then(r => r.data);
