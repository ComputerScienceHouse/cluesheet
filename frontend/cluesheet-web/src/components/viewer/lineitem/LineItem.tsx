import { Clue, UserClue } from "@/lib/types";

import styles from "./page.module.scss";
import { useState } from "react";
interface LineItemProps {
  parent: Clue | UserClue | null;
  clue: Clue | UserClue;
}

export default function LineItem({ parent, clue }: LineItemProps) {
  // Initialize state to track if the checkbox is checked
  const [isChecked, setIsChecked] = useState(
    (clue as UserClue).completions !== undefined &&
      Number((clue as UserClue).completions) > 0,
  );

  // Handle checkbox change
  const handleClueChecked = () => {
    setIsChecked((prevChecked) => !prevChecked);
    (clue as UserClue).completions = 1;
    // You can also perform any additional logic here, like updating the clue state
  };

  const shouldLock = (parent !== null && (parent as UserClue).completions <= 0);

  let tags = [];
  if (parent != null) {
    tags = clue.tags.filter(item => !parent.tags.includes(item));
  } else {
    tags = clue.tags;
  }

  return (
    <>
      <div className={`${styles.lineItem} ${shouldLock ? styles.locked : ''}`}>
        {(clue as UserClue).completions !== undefined && (
          <input
            type="checkbox"
            className={styles.customCheckboxInput}
            onChange={handleClueChecked}
            checked={isChecked}
          />
        )}
        <h5 className={styles.points}>
          {clue.points}:
        </h5>
        <h5 className={styles.description}>
          {clue.description}
        </h5>
        <div className={styles.tags}>
          {tags.length > 0 &&
          tags.map((tag: string, index) => (
          <p className={styles.tag} key={index}>{tag}</p>
          ))}
        </div>
      </div>
    </>
  );
}
