import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import ClueGroup from "@/components/viewer/cluegroup/ClueGroup";
import { mockCluesheet } from "../../../../tests/lib/data";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetViewer({
  params,
}: {
  params: Promise<{ cluesheet_id: string }>;
}) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);
  
  const cluesheet = mockCluesheet;

  return (
    <>
      <main>
        <Container>
          <div className={styles.header}>
            <h1>{cluesheet.title}</h1>
          </div>
          <div className={styles.clueList}>
            {cluesheet.clues.length > 0 &&
              cluesheet.clues.map((clue: Clue, index) => (
                <ClueGroup key={`${clue.id}`} clue={clue} />
              ))}
          </div>
        </Container>
      </main>
    </>
  );
}
