"use client"

import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import { ClueList } from "@/components/viewer/ClueList/ClueList";
import { mockUserCluesheet } from "../../../../tests/lib/data";

interface CluesheetFormProps {
		cluesheet_id: string;
		keycloak_uid: string;
}

export default async function CluesheetForm({cluesheet_id, keycloak_uid}: CluesheetFormProps) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}/user/${keycloak_uid}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  const cluesheet = mockUserCluesheet;

  return (
    <>
          <div className={styles.header}>
            <h1>{cluesheet.title}</h1>
            <div className={styles.pointCounter}>
              <PointCounter points={cluesheet.user_points} />
              <h3>Points</h3>
            </div>
          </div>
          <div className={styles.clueList}>
            {cluesheet.clues.length > 0 && (
              <ClueList ancestors={[]} clues={cluesheet.clues} edit={true} />
            )}
          </div>
    </>
  );
}
