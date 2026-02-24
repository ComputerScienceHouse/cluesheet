import Image from "next/image";
import styles from "./page.module.scss";

export default function Home() {
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <p>Welcome to willard's world</p>
      </main>
      <footer className={styles.footer}></footer>
    </div>
  );
}
