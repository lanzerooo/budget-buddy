import { useState, type FormEvent } from "react";
import styles from "./LoginRegister.module.scss";
import axios, { AxiosError } from "axios";

function LoginRegister() {
    const [formType, setFormType] = useState<"login" | "register">("login");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");

    const handleFormSwitch = (type: "login" | "register") => {
        setFormType(type);
        setEmail("");
        setPassword("");
        setName("");
        setConfirmPassword("");
        setError("");
    };

    async function handleLogin(email: string, password: string): Promise<void> {
        try {
            const response = await axios.post<{ token: string }>('http://localhost:8080/login', {
                email,
                password,
            });
            localStorage.setItem('token', response.data.token);
            console.log('Успешный вход:', response.data.token);
            setError("");
        } catch (error: unknown) {
            if (axios.isAxiosError(error)) {
                const axiosError = error as AxiosError<{ message?: string }>;
                setError(axiosError.response?.data?.message || 'Что-то пошло не так');
            } else {
                setError('Неизвестная ошибка');
            }
        }
    }

    async function handleRegister(email: string, password: string, name: string): Promise<void> {
        if (password !== confirmPassword) {
            setError('Пароли не совпадают');
            return;
        }
        try {
            const response = await axios.post<{ token: string }>('http://localhost:8080/register', {
                email,
                password,
                name,
            });
            localStorage.setItem('token', response.data.token);
            console.log('Успешная регистрация:', response.data.token);
            setError("");
        } catch (error: unknown) {
            if (axios.isAxiosError(error)) {
                const axiosError = error as AxiosError<{ message?: string }>;
                setError(axiosError.response?.data?.message || 'Что-то пошло не так');
            } else {
                setError('Неизвестная ошибка');
            }
        }
    }

    const handleLoginSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        handleLogin(email, password);
    };

    const handleRegisterSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        handleRegister(email, password, name);
    };

    return (
        <div className={styles.container}>
            <div className={styles.tab}>
                <button
                    className={`${styles.tabButton} ${formType === "login" ? styles.active : ""}`}
                    onClick={() => handleFormSwitch("login")}
                >
                    Логин
                </button>
                <button
                    className={`${styles.tabButton} ${formType === "register" ? styles.active : ""}`}
                    onClick={() => handleFormSwitch("register")}
                >
                    Регистрация
                </button>
            </div>
            {error && <p className={styles.error}>{error}</p>}
            {formType === "login" ? (
                <form onSubmit={handleLoginSubmit} className={styles.form}>
                    <h2>Вход</h2>
                    <input
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <button type="submit" className={styles.button}>Войти</button>
                </form>
            ) : (
                <form onSubmit={handleRegisterSubmit} className={styles.form}>
                    <h2>Регистрация</h2>
                    <input
                        type="text"
                        placeholder="Имя"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <input
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <input
                        type="password"
                        placeholder="Подтвердите пароль"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                        className={styles.input}
                    />
                    <button type="submit" className={styles.button}>Зарегистрироваться</button>
                </form>
            )}
        </div>
    );
}

export default LoginRegister;