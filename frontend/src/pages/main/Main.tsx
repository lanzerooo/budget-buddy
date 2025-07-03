import { Footer } from '../../components/footer/component'
import { Header } from '../../components/header/component'
import styles from './Main.module.scss'

function Main() {
 

  return (
    <>
    <Header/>
      <div className={styles.main_section}>
        <h1 className={styles.title}>Welcome to BudgetBuddy</h1>
        <section className={styles.income}>
            
        </section>
      </div>
      <Footer/>
    </>
  )
}

export default Main
