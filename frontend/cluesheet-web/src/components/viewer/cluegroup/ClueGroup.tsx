"use client";
import { Clue, UserClue } from "@/lib/types";
import styles from "./page.module.scss";
import LineItem from "../lineitem/LineItem";

interface ClueGroupProps {
  clue: Clue | UserClue;
}

export default function ClueGroup({ clue }: ClueGroupProps) {
  return (
    <>
      <div className={`${styles.clueBody} ${/*possible ? '' : styles.locked*/''}`}>
        <LineItem clue={clue}/>
        <div className={styles.children}>
          {clue.children.length > 0 &&
            clue.children.map((clue: Clue, index) => (
              <ClueGroup key={`${clue.id}`} clue={clue}/>
            ))}
        </div>
      </div>
    </>
  );
}
