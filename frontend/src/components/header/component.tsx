import React from 'react';
import { Link } from 'react-router-dom';
import styles from './Header.module.scss';

// Интерфейс для пропсов (опционально, для расширяемости)
interface HeaderProps {
  appName?: string; // Название приложения для отображения в шапке
}

// Создаем компонент Header
export const Header: React.FC<HeaderProps> = ({ appName = 'BudgetBuddy' }) => {
  return (
    <header className={styles.header}>
      <div className={styles.container}>
        {/* Логотип или название приложения */}
        <h1 className={styles.logo}>{appName}</h1>
        
        {/* Навигация */}
        <nav className={styles.nav}>
          <Link to="/main" className={styles.navLink}>Главная</Link>
          <Link to="/transactions" className={styles.navLink}>Транзакции</Link>
          <Link to="/goals" className={styles.navLink}>Цели</Link>
          <Link to="/analytics" className={styles.navLink}>Аналитика</Link>
        </nav>
      </div>
    </header>
  );
};