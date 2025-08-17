export interface User {
    id: string;
    username: string;
    email: string;
};

export interface AuthState {
    user: User | null;
    accessToken: string | null;
    isAuthenticated: boolean;
};
