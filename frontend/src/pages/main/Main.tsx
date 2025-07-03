import { useState, useEffect } from 'react'; // Добавляем useEffect для автоматического вызова
import { Footer } from '../../components/footer/component';
import { Header } from '../../components/header/component';
import styles from './Main.module.scss';

interface Transaction {
  id: number;
  type: 'income' | 'expense'; // Тип транзакции
  amount: number;              // Сумма
  description: string;         // Описание
  date: string;                // Дата в формате ISO (например, "2025-07-03")
  category_id?: number;        // Необязательное поле категории
}

function Main() {
  const API_BASE_URL = 'http://localhost'; // Базовый URL, порт добавим динамически
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Функция для получения транзакций
  async function handleTransaction() {
    const token = localStorage.getItem('token');
    if (!token) {
      setError('Токен не найден. Пожалуйста, войдите снова.');
      setLoading(false);
      return;
    }

    const transactionsUrl = `${API_BASE_URL}:${import.meta.env.VITE_FINANCE_SERVICE_PORT}/transactions`;
    try {
      const response = await fetch(transactionsUrl, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Ошибка: ${await response.text()}`);
      }

      const data = await response.json();
      setTransactions(data);
    } catch (err) {
      setError((err as Error).message || 'Неизвестная ошибка');
    } finally {
      setLoading(false);
    }
  }

  // Автоматический вызов при монтировании компонента
  useEffect(() => {
    handleTransaction();
  }, []); // Пустой массив зависимостей — вызывается один раз при загрузке

  return (
    <>
      <Header />
      <div className={styles.main_section}>
        <h1 className={styles.title}>Welcome to BudgetBuddy</h1>
        <section className={styles.income}>
          {loading && <p>Загрузка...</p>}
          {error && <p style={{ color: 'red' }}>Ошибка: {error}</p>}
          {!loading && !error && transactions.length > 0 && (
            <ul>
              {transactions.map((tx) => (
                <li key={tx.id}>
                  {tx.date} - {tx.type === 'income' ? 'Доход' : 'Расход'}: {tx.amount} ₽ - {tx.description}
                </li>
              ))}
            </ul>
          )}
          {!loading && !error && transactions.length === 0 && <p>Нет транзакций.</p>}
        </section>
      </div>
      <Footer />
    </>
  );
}

export default Main;