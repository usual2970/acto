

export type Response<T> = {
    code: number;
    message: string;
    data: T;
};


export type User = {
    username: string;
    password: string;
};