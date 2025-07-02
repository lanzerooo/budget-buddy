import { useState } from "react";
import styles from "./LoginRegister.module.scss";

function LoginRegister() {
    const [formType, setFormType] = useState<"login" | "register">("login");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const handleFormSwitch = (type: "login" | "register") => {
        setFormType(type);
        setEmail("");
        setPassword("");
        setName("");
        setConfirmPassword("");
    };

    const handleLoginSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log("Login:", { email, password });
        // Здесь можно добавить логику отправки данных на сервер
    };

    const handleRegisterSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            console.error("Пароли не совпадают");
            return;
        }
        console.log("Register:", { name, email, password });
        // Здесь можно добавить логику отправки данных на сервер
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
            {formType === "login" ? (
                <form onSubmit={handleLoginSubmit} className={styles.form}>
                    <h2>Вход</h2>
                    <input
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                    <button type="submit">Войти</button>
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
                    />
                    <input
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Подтвердите пароль"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                    />
                    <button type="submit">Зарегистрироваться</button>
                </form>
            )}
        </div>
    );
}

export default LoginRegister;