import { Clue, UserClue } from "@/lib/types";

import styles from "./page.module.scss";
import { useState } from "react";
interface LineItemProps {
  clue: Clue | UserClue;
}

export default function LineItem({ clue }: LineItemProps) {
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

  return (
    <>
      <div className={styles.lineItem}>
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
          {clue.tags.length > 0 &&
          clue.tags.map((tag: string, index) => (
          <p className={styles.tag} key={index}>{tag}</p>
          ))}
        </div>
      </div>
    </>
  );
}
