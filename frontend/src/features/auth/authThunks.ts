import { createAsyncThunk } from '@reduxjs/toolkit';
import { loginSuccess, logout } from './authSlice';
import { User } from '@/types/auth';
import { RootState } from '@/app/store';

export const login = createAsyncThunk(
    'auth/login',
    async(credentials: { username: string; password: string }, { dispatch, rejectWithValue}) => {
        try {
            const response = await axios.post('/api/auth/login', credentials);

            const rawUser = response.data.user;
            const sanitizedUser: User = {
                id: String(rawUser.id),
                username: rawUser.username,
                email: rawUser.email,
            };

            const accessToken = response.data.access_token;

            dispatch(loginSuccess({ user: sanitizedUser, accessToken }));
        } catch (error: any) {
            return rejectWithValue(error.response?.data || 'Login failed');
        }
    }
);

export const logoutUser = createAsyncThunk('auth/logoutUser', async(_, { dispatch }) => {
    dispatch(logout());
});
