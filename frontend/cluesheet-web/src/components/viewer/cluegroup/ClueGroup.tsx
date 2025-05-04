"use client";
import { Clue } from "@/lib/types";
import styles from "./page.module.scss";
import LineItem from "../lineitem/LineItem";

interface ClueGroupProps {
  clue: Clue;
}

export default function ClueGroup({ clue }: ClueGroupProps) {
  return (
    <>
      <div className={styles.clueBody}>
        <LineItem clue={clue}/>
        <div className={styles.children}>
          {clue.children.length > 0 &&
            clue.children.map((clue: Clue, index) => (
              <ClueGroup key={`${clue.id}-${index}`} clue={clue} />
            ))}
        </div>
      </div>
    </>
  );
}
