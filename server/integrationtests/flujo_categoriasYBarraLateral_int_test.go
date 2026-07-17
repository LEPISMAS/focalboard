package integrationtests

import (
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/require"
)

func TestINT0601CrearCategoriaPersistenciaStore(t *testing.T) {
	t.Log("INT-06-01: Crear una categoría vía API y verificar persistencia en Store")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	categoryName := "Custom Sidebar Category"
	categoryToCreate := model.Category{
		Name:   categoryName,
		TeamID: testTeamID,
		UserID: user1.ID,
	}

	createdCat, resp := th.Client.CreateCategory(categoryToCreate)
	th.CheckOK(resp)
	require.NotNil(t, createdCat)
	require.NotEmpty(t, createdCat.ID)
	require.Equal(t, categoryName, createdCat.Name)

	// Verificar persistencia en base de datos (Store)
	persistedCat, err := th.Server.Store().GetCategory(createdCat.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedCat)
	require.Equal(t, categoryName, persistedCat.Name)
	require.Equal(t, testTeamID, persistedCat.TeamID)
	require.Equal(t, user1.ID, persistedCat.UserID)
}

func TestINT0602AsociarTableroACategoria(t *testing.T) {
	t.Log("INT-06-02: Mover un tablero a una categoría personalizada y verificar la asociación en el Store")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)

	// Crear categoría personalizada
	createdCat, resp := th.Client.CreateCategory(model.Category{
		Name:   "Custom Category 2",
		TeamID: testTeamID,
		UserID: user1.ID,
	})
	th.CheckOK(resp)
	require.NotNil(t, createdCat)

	// Mover tablero a la nueva categoría
	updateResp := th.Client.UpdateCategoryBoard(testTeamID, createdCat.ID, board.ID)
	th.CheckOK(updateResp)

	// Verificar asociación directamente en la base de datos (Store)
	categoryBoards, err := th.Server.Store().GetUserCategoryBoards(user1.ID, testTeamID)
	require.NoError(t, err)

	var foundCategory *model.CategoryBoards
	for _, cb := range categoryBoards {
		if cb.ID == createdCat.ID {
			foundCategory = &cb
			break
		}
	}
	require.NotNil(t, foundCategory, "Category not found in store")
	require.Len(t, foundCategory.BoardMetadata, 1)
	require.Equal(t, board.ID, foundCategory.BoardMetadata[0].BoardID)
}

func TestINT0603ObtenerCategoriasConTableros(t *testing.T) {
	t.Log("INT-06-03: Obtener categorías del sidebar y verificar que incluyan los tableros asociados")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board1 := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	board2 := th.CreateBoard(testTeamID, model.BoardTypeOpen)

	// Crear categoría personalizada
	createdCat, resp := th.Client.CreateCategory(model.Category{
		Name:   "Sidebar Boards Category",
		TeamID: testTeamID,
		UserID: user1.ID,
	})
	th.CheckOK(resp)
	require.NotNil(t, createdCat)

	// Asociar ambos tableros a la categoría personalizada
	updateResp1 := th.Client.UpdateCategoryBoard(testTeamID, createdCat.ID, board1.ID)
	th.CheckOK(updateResp1)
	updateResp2 := th.Client.UpdateCategoryBoard(testTeamID, createdCat.ID, board2.ID)
	th.CheckOK(updateResp2)

	// Obtener categorías desde el sidebar vía API
	categoriesWithBoards, apiResp := th.Client.GetUserCategoryBoards(testTeamID)
	th.CheckOK(apiResp)
	require.NotEmpty(t, categoriesWithBoards)

	var foundCategory *model.CategoryBoards
	for _, cb := range categoriesWithBoards {
		if cb.ID == createdCat.ID {
			foundCategory = &cb
			break
		}
	}
	require.NotNil(t, foundCategory)
	require.Len(t, foundCategory.BoardMetadata, 2)

	// Verificar que ambos IDs de tableros se encuentren en BoardMetadata
	var foundBoard1, foundBoard2 bool
	for _, bm := range foundCategory.BoardMetadata {
		if bm.BoardID == board1.ID {
			foundBoard1 = true
		}
		if bm.BoardID == board2.ID {
			foundBoard2 = true
		}
	}
	require.True(t, foundBoard1)
	require.True(t, foundBoard2)
}

func TestINT0604EliminarCategoriaReasignarDefault(t *testing.T) {
	t.Log("INT-06-04: Eliminar una categoría y verificar que los tableros se muevan a la categoría por defecto")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)

	// Crear categoría personalizada
	createdCat, resp := th.Client.CreateCategory(model.Category{
		Name:   "Category to Delete",
		TeamID: testTeamID,
		UserID: user1.ID,
	})
	th.CheckOK(resp)
	require.NotNil(t, createdCat)

	// Mover tablero a esta categoría
	updateResp := th.Client.UpdateCategoryBoard(testTeamID, createdCat.ID, board.ID)
	th.CheckOK(updateResp)

	// Eliminar categoría personalizada
	deleteResp := th.Client.DeleteCategory(testTeamID, createdCat.ID)
	th.CheckOK(deleteResp)

	// Obtener categorías boards actualizados
	categoriesWithBoards, apiResp := th.Client.GetUserCategoryBoards(testTeamID)
	th.CheckOK(apiResp)

	// La categoría creada debe haber sido eliminada
	for _, cb := range categoriesWithBoards {
		require.NotEqual(t, createdCat.ID, cb.ID)
	}

	// El tablero debe haber sido reasignado a la categoría por defecto (system: Boards)
	var defaultCategory *model.CategoryBoards
	for _, cb := range categoriesWithBoards {
		if cb.Type == model.CategoryTypeSystem && cb.Name == "Boards" {
			defaultCategory = &cb
			break
		}
	}
	require.NotNil(t, defaultCategory, "Default Boards category not found")

	var boardFoundInDefault bool
	for _, bm := range defaultCategory.BoardMetadata {
		if bm.BoardID == board.ID {
			boardFoundInDefault = true
			break
		}
	}
	require.True(t, boardFoundInDefault, "Board was not moved back to the default category")
}

func TestINT0605ReordenarCategoriasPersistencia(t *testing.T) {
	t.Log("INT-06-05: Reordenar categorías y verificar que el nuevo orden se persista")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()

	// Crear dos categorías
	cat1, resp1 := th.Client.CreateCategory(model.Category{
		Name:   "A Category 1",
		TeamID: testTeamID,
		UserID: user1.ID,
	})
	th.CheckOK(resp1)

	cat2, resp2 := th.Client.CreateCategory(model.Category{
		Name:   "B Category 2",
		TeamID: testTeamID,
		UserID: user1.ID,
	})
	th.CheckOK(resp2)

	// Obtener categorías boards inicialmente
	initialCats, apiResp := th.Client.GetUserCategoryBoards(testTeamID)
	th.CheckOK(apiResp)
	require.NotEmpty(t, initialCats)

	// Construir el nuevo orden incluyendo TODOS los IDs de categorías.
	// El backend exige que la lista tenga el mismo número de entradas que las
	// categorías existentes en la base de datos (incluyendo la categoría de
	// sistema "Boards"). Ponemos cat2 antes que cat1 y al final las demás.
	customIDs := map[string]bool{cat1.ID: true, cat2.ID: true}
	newOrder := []string{cat2.ID, cat1.ID}
	for _, cb := range initialCats {
		if !customIDs[cb.ID] {
			newOrder = append(newOrder, cb.ID)
		}
	}

	// Reordenar categorías
	reorderedIDs, reorderResp := th.Client.ReorderCategories(testTeamID, newOrder)
	th.CheckOK(reorderResp)
	require.NotEmpty(t, reorderedIDs)

	// Volver a obtener las categorías para verificar la persistencia
	updatedCats, getResp := th.Client.GetUserCategoryBoards(testTeamID)
	th.CheckOK(getResp)

	// Buscar las posiciones de cat1 y cat2 en las categorías actualizadas
	var indexCat1, indexCat2 int = -1, -1
	for idx, cb := range updatedCats {
		if cb.ID == cat1.ID {
			indexCat1 = idx
		} else if cb.ID == cat2.ID {
			indexCat2 = idx
		}
	}

	require.NotEqual(t, -1, indexCat1)
	require.NotEqual(t, -1, indexCat2)

	// Como cat2 está antes que cat1 en el nuevo orden, indexCat2 debe ser menor que indexCat1
	require.Less(t, indexCat2, indexCat1, "cat2 should be ordered before cat1")
}
