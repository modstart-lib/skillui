import axios from 'axios';

const http = axios.create({
  baseURL: 'https://skillui.com/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
});

export function setApiBaseUrl(url: string) {
  http.defaults.baseURL = url;
}

export default http;
