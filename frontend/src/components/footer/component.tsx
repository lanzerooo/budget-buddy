import React from 'react';
import { Link } from 'react-router-dom'; // Используем Link вместо a для навигации
import styles from './Footer.module.scss';

// Интерфейс для пропсов (опционально, для расширяемости)
interface FooterProps {
  companyName?: string; // Название компании для копирайта
}

// Создаем компонент Footer
export const Footer: React.FC<FooterProps> = ({ companyName = 'MyApp' }) => {
  return (
    <footer className={styles.footer}>
      <div className={styles.container}>
        <p className={styles.copyright}>
          © {new Date().getFullYear()} {companyName}. Все права защищены.
        </p>
        <nav className={styles.nav}>
          <Link to="/about" className={styles.navLink}>О нас</Link>
          <Link to="/contact" className={styles.navLink}>Контакты</Link>
          <Link to="/privacy" className={styles.navLink}>Политика конфиденциальности</Link>
        </nav>
      </div>
    </footer>
  );
};