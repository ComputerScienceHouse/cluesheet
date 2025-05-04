"use client"
import { Clue } from "@/lib/types";
import styles from "./page.module.scss";

interface ClueGroupProps {
  clue: Clue;
}

export default function ClueGroup({ clue }: ClueGroupProps) {
  function handleClueChecked() {
    console.log("Checked!");
  }

  return (
    <>
      <div className={styles.clueBody}>
        <div className={styles.lineItem}>
          <input type="checkbox" className={styles.checkbox} onClick={handleClueChecked} checked={Number(clue.completions) > 0 ? true : false}/>
          <h3>{clue.points}: {clue.description}</h3>
        </div>
        <div className={styles.children}>
        {clue.children.length > 0 &&
          clue.children.map((clue: Clue, index) => (
            <ClueGroup key={clue.id} clue={clue} />
          ))}
          </div>
          </div>
    </>
  );
}
