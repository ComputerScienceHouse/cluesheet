"use client";
import { Clue, UserClue } from "@/lib/types";
import styles from "./page.module.scss";
import LineItem from "../lineitem/LineItem";

interface ClueGroupProps {
  parent: Clue | UserClue | null;
  clue: Clue | UserClue;
}

export default function ClueGroup({ parent, clue }: ClueGroupProps) {
  if (parent !== null) {
    console.log(`OG Tags: ${parent.tags}`);
    parent.tags = parent?.tags.concat(clue.tags);
    console.log(`New Tags: ${parent.tags}`);
  } else {
    console.log(`Error: parent is ${parent}`);
  }

  return (
    <>
      <div className={`${styles.clueBody} ${/*possible ? '' : styles.locked*/''}`}>
        <LineItem clue={clue} parent={parent}/>
        <div className={styles.children}>
          {clue.children.length > 0 &&
            clue.children.map((child: Clue, index) => (
              <ClueGroup key={`${child.id}`} clue={child} parent={clue}/>
            ))}
        </div>
      </div>
    </>
  );
}
