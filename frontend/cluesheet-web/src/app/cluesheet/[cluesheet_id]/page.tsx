import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import { mockCluesheet } from "../../../../tests/lib/data";
import { ClueList } from "@/components/viewer/ClueList/ClueList";
import CluesheetForm from "@/components/form/CluesheetForm/CluesheetForm";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetViewer({
  params,
}: {
  params: Promise<{ cluesheet_id: string }>;
}) {
  //const cluesheet = mockCluesheet.clues;
  //const cluesheet: Array<Clue> = [];
  //const { cluesheet_id } = await params;
  //const cluesheet = mockUserCluesheet;

  return (
    <>
      <main>
        <Container>
          <CluesheetForm
            cluesheet_id={"e"}
            keycloak_uid={"chom"}
            edit={false}
          />
        </Container>
      </main>
    </>
  );
}
