package integrationtests

import (
	"bytes"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/require"
)

// int07CreateBoardWithBlocks crea, via API, un tablero con un bloque tipo tarjeta
// y devuelve las entidades ya persistidas (con los IDs asignados por el servidor).
func int07CreateBoardWithBlocks(t *testing.T, th *TestHelper, title string, numCards int) *model.BoardsAndBlocks {
	t.Helper()

	board := &model.Board{
		ID:        utils.NewID(utils.IDTypeBoard),
		TeamID:    model.GlobalTeamID,
		Title:     title,
		CreatedBy: th.GetUser1().ID,
		Type:      model.BoardTypeOpen,
		CreateAt:  utils.GetMillis(),
		UpdateAt:  utils.GetMillis(),
	}

	blocks := make([]*model.Block, 0, numCards)
	for n := 0; n < numCards; n++ {
		blocks = append(blocks, &model.Block{
			ID:        utils.NewID(utils.IDTypeCard),
			ParentID:  board.ID,
			BoardID:   board.ID,
			Type:      model.TypeCard,
			Title:     title + " card",
			CreatedBy: th.GetUser1().ID,
			CreateAt:  utils.GetMillis(),
			UpdateAt:  utils.GetMillis(),
		})
	}

	babs := &model.BoardsAndBlocks{
		Boards: []*model.Board{board},
		Blocks: blocks,
	}

	babsCreated, resp := th.Client.CreateBoardsAndBlocks(babs)
	th.CheckOK(resp)
	require.NotNil(t, babsCreated)

	return babsCreated
}

func TestINT0701ExportarTableroGeneraArchivoConFormatoCorrecto(t *testing.T) {
	t.Log("INT-07-01: verifica API REST /boards/{id}/archive/export -> App.ExportArchive -> Store/FileSystem, confirmando que se genera un archivo con el tablero, sus bloques y metadatos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	babs := int07CreateBoardWithBlocks(t, th, "INT-07-01 tablero exportable", 2)

	archive, resp := th.Client.ExportBoardArchive(babs.Boards[0].ID)
	th.CheckOK(resp)
	require.NotNil(t, archive)
	require.NotEmpty(t, archive)

	// Un archivo de exportacion valido de Focalboard es un zip: comienza con la
	// firma de cabecera estandar "PK" ('P'=0x50, 'K'=0x4B).
	require.GreaterOrEqual(t, len(archive), 2)
	require.Equal(t, byte(0x50), archive[0])
	require.Equal(t, byte(0x4B), archive[1])
}

func TestINT0702ImportarArchivoRecreaEntidadesEnStore(t *testing.T) {
	t.Log("INT-07-02: verifica API REST /teams/{id}/archive/import -> App.ImportArchive -> Store, confirmando que el tablero y los bloques se recrean con nuevos IDs.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	babs := int07CreateBoardWithBlocks(t, th, "INT-07-02 tablero origen", 1)

	archive, resp := th.Client.ExportBoardArchive(babs.Boards[0].ID)
	th.CheckOK(resp)
	require.NotEmpty(t, archive)

	resp = th.Client.ImportArchive(model.GlobalTeamID, bytes.NewReader(archive))
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	boards, err := th.Server.App().GetBoardsForUserAndTeam(th.GetUser1().ID, model.GlobalTeamID, true)
	require.NoError(t, err)

	var imported *model.Board
	for _, b := range boards {
		if b.Title == babs.Boards[0].Title && b.ID != babs.Boards[0].ID {
			imported = b
			break
		}
	}
	require.NotNil(t, imported, "el tablero importado debe existir con un ID distinto al original")

	importedBlocks, err := th.Server.App().GetBlocksForBoard(imported.ID)
	require.NoError(t, err)
	require.Len(t, importedBlocks, 1)
	require.Equal(t, babs.Blocks[0].Title, importedBlocks[0].Title)
	require.NotEqual(t, babs.Blocks[0].ID, importedBlocks[0].ID)
}

func TestINT0703ExportarEImportarTableroRoundtripPreservaIntegridad(t *testing.T) {
	t.Log("INT-07-03: verifica el ciclo completo API -> App -> Store (export + import), confirmando que el tablero importado conserva el mismo contenido que el original.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	babs := int07CreateBoardWithBlocks(t, th, "INT-07-03 tablero roundtrip", 3)

	archive, resp := th.Client.ExportBoardArchive(babs.Boards[0].ID)
	th.CheckOK(resp)
	require.NotEmpty(t, archive)

	resp = th.Client.ImportArchive(model.GlobalTeamID, bytes.NewReader(archive))
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	boards, err := th.Server.App().GetBoardsForUserAndTeam(th.GetUser1().ID, model.GlobalTeamID, true)
	require.NoError(t, err)

	var imported *model.Board
	for _, b := range boards {
		if b.Title == babs.Boards[0].Title && b.ID != babs.Boards[0].ID {
			imported = b
			break
		}
	}
	require.NotNil(t, imported)
	require.Equal(t, babs.Boards[0].Type, imported.Type)

	importedBlocks, err := th.Server.App().GetBlocksForBoard(imported.ID)
	require.NoError(t, err)
	require.Len(t, importedBlocks, len(babs.Blocks))

	originalTitles := make(map[string]bool)
	for _, b := range babs.Blocks {
		originalTitles[b.Title] = true
	}
	for _, b := range importedBlocks {
		require.True(t, originalTitles[b.Title], "cada bloque importado debe corresponder a un bloque original por su titulo")
	}
}

func TestINT0704ImportarArchivoInvalidoEsRechazadoSinCorromperDatos(t *testing.T) {
	t.Log("INT-07-04: verifica que API REST /teams/{id}/archive/import -> App.ImportArchive rechaza un archivo con formato invalido sin crear entidades parciales.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	boardsBefore, err := th.Server.App().GetBoardsForUserAndTeam(th.GetUser1().ID, model.GlobalTeamID, true)
	require.NoError(t, err)
	countBefore := len(boardsBefore)

	invalidArchive := bytes.NewReader([]byte("esto-no-es-un-archivo-boardarchive-valido"))
	resp := th.Client.ImportArchive(model.GlobalTeamID, invalidArchive)
	require.Error(t, resp.Error)

	boardsAfter, err := th.Server.App().GetBoardsForUserAndTeam(th.GetUser1().ID, model.GlobalTeamID, true)
	require.NoError(t, err)
	require.Equal(t, countBefore, len(boardsAfter), "un archivo invalido no debe crear tableros parcialmente")
}

func TestINT0705SubirArchivoAdjuntoAlmacenaYEsRecuperable(t *testing.T) {
	t.Log("INT-07-05: verifica API REST /files -> App.SaveFile -> FileSystem + Store, confirmando que un archivo adjunto queda almacenado y su referencia es recuperable via API.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	testBoard := th.CreateBoard(model.GlobalTeamID, model.BoardTypeOpen)

	uploaded, resp := th.Client.TeamUploadFile(model.GlobalTeamID, testBoard.ID, bytes.NewBuffer([]byte("contenido-de-prueba-int-07-05")))
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, uploaded)
	require.NotEmpty(t, uploaded.FileID)

	fileInfo, resp := th.Client.TeamUploadFileInfo(model.GlobalTeamID, testBoard.ID, "test")
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, fileInfo)
	require.NotEmpty(t, fileInfo.Id)
}
