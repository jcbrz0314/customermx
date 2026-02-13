package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

// Configuration
const (
	ExcelFile    = "eventos.xlsx"
	SheetName    = "General"
	OutputDir    = "backend/migrations/import"
	DBConnString = "postgres://postgres:postgres@localhost:5432/customermx?sslmode=disable"

	// Column mappings (1-indexed)
	ColMarcaEvento      = 2  // B
	ColTipoEvento       = 3  // C
	ColOrganizador      = 4  // D
	ColNombreEvento     = 5  // E
	ColFecha            = 6  // F
	ColAno              = 7  // G
	ColDiasEvento       = 8  // H
	ColEstado           = 9  // I
	ColCiudad           = 10 // J
	ColSede             = 11 // K
	ColDistribuidor     = 12 // L
	ColEdecanes         = 13 // M
	ColMontaje          = 14 // N
	ColPromocional      = 15 // O
	ColAsistencia       = 16 // P
	ColDinamicas        = 17 // Q
	ColDatosLevantados  = 18 // R
	ColProspectos       = 19 // S
	ColVehiculosStart   = 20 // T (inicio de vehículos)
	ColCalificacion     = 60 // BH
	ColComentarios      = 61 // BI

	// Data starts at row 3
	DataStartRow = 3
)

// Domain models
type Brand struct {
	ID   uuid.UUID
	Name string
}

type Vehicle struct {
	ID        uuid.UUID
	BrandID   uuid.UUID
	ModelName string
}

type Event struct {
	ID           uuid.UUID
	BrandID      uuid.UUID
	EventType    string
	Organizer    string
	Name         string
	StartDate    time.Time
	Year         int
	DurationDays int
	State        string
	City         string
	Venue        string
	Dealer       string
	Status       string
}

type EventReport struct {
	ID              uuid.UUID
	EventID         uuid.UUID
	HostessCount    *int
	SetupVendor     *string
	HasPromotional  *bool
	Attendees       *int
	ActivitiesCount *int
	LeadsCollected  *int
	Prospects       *int
	DealerRating    *int
	Comments        *string
	Completed       bool
}

type EventVehicle struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	VehicleID uuid.UUID
	Quantity  int
}

type VehicleColumn struct {
	ColIndex  int
	BrandName string
	ModelName string
}

// In-memory catalog
type Catalog struct {
	Brands            map[string]*Brand            // name -> Brand
	Vehicles          map[string]*Vehicle          // "BrandName:ModelName" -> Vehicle
	VehiclesByBrandID map[uuid.UUID][]*Vehicle     // brand_id -> []Vehicle
	NewBrands         []*Brand                     // brands to insert
	NewVehicles       []*Vehicle                   // vehicles to insert
	Events            []*Event                     // events to insert
	EventReports      []*EventReport               // reports to insert
	EventVehicles     []*EventVehicle              // event_vehicles to insert
	VehicleColumns    map[int]*VehicleColumn       // colIndex -> vehicle info
}

func main() {
	log.Println("🚀 Iniciando importación de eventos desde Excel...")

	// Step 1: Open Excel file
	log.Printf("📂 Abriendo archivo Excel: %s", ExcelFile)
	f, err := excelize.OpenFile(ExcelFile)
	if err != nil {
		log.Fatalf("❌ Error al abrir archivo Excel: %v", err)
	}
	defer f.Close()

	// Step 2: Initialize catalog from database
	log.Println("🗄️  Conectando a base de datos...")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, DBConnString)
	if err != nil {
		log.Fatalf("❌ Error al conectar a base de datos: %v", err)
	}
	defer pool.Close()

	catalog := &Catalog{
		Brands:            make(map[string]*Brand),
		Vehicles:          make(map[string]*Vehicle),
		VehiclesByBrandID: make(map[uuid.UUID][]*Vehicle),
		VehicleColumns:    make(map[int]*VehicleColumn),
	}

	log.Println("📚 Cargando catálogo de marcas y vehículos existentes...")
	if err := loadExistingCatalog(ctx, pool, catalog); err != nil {
		log.Fatalf("❌ Error al cargar catálogo: %v", err)
	}
	log.Printf("✅ Catálogo cargado: %d marcas, %d vehículos", len(catalog.Brands), len(catalog.Vehicles))

	// Step 3: Parse Excel headers (row 1 and 2)
	log.Println("📋 Parseando encabezados de Excel...")
	if err := parseVehicleHeaders(f, catalog); err != nil {
		log.Fatalf("❌ Error al parsear encabezados: %v", err)
	}
	log.Printf("✅ Encabezados parseados: %d columnas de vehículos", len(catalog.VehicleColumns))

	// Step 4: Parse data rows (row 3+)
	log.Println("📊 Parseando datos de eventos...")
	if err := parseEventData(f, catalog); err != nil {
		log.Fatalf("❌ Error al parsear datos: %v", err)
	}
	log.Printf("✅ Datos parseados: %d eventos, %d reportes, %d event_vehicles",
		len(catalog.Events), len(catalog.EventReports), len(catalog.EventVehicles))

	// Step 5: Show warnings for new items
	if len(catalog.NewBrands) > 0 {
		log.Printf("⚠️  Se encontraron %d marcas nuevas que se insertarán:", len(catalog.NewBrands))
		for _, b := range catalog.NewBrands {
			log.Printf("   - %s", b.Name)
		}
	}
	if len(catalog.NewVehicles) > 0 {
		log.Printf("⚠️  Se encontraron %d vehículos nuevos que se insertarán:", len(catalog.NewVehicles))
		for _, v := range catalog.NewVehicles {
			brandName := getBrandNameByID(catalog, v.BrandID)
			log.Printf("   - %s: %s", brandName, v.ModelName)
		}
	}

	// Step 6: Generate SQL files
	log.Println("📝 Generando archivos SQL...")
	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		log.Fatalf("❌ Error al crear directorio de salida: %v", err)
	}
	if err := generateSQLFiles(catalog); err != nil {
		log.Fatalf("❌ Error al generar archivos SQL: %v", err)
	}
	log.Println("✅ Archivos SQL generados en:", OutputDir)

	// Step 7: Execute SQL in database
	log.Println("💾 Ejecutando SQL en base de datos...")
	if err := executeSQL(ctx, pool, catalog); err != nil {
		log.Fatalf("❌ Error al ejecutar SQL: %v", err)
	}

	log.Println("✅ ¡Importación completada exitosamente!")
	log.Printf("📊 Resumen:")
	log.Printf("   - Marcas nuevas: %d", len(catalog.NewBrands))
	log.Printf("   - Vehículos nuevos: %d", len(catalog.NewVehicles))
	log.Printf("   - Eventos insertados: %d", len(catalog.Events))
	log.Printf("   - Reportes insertados: %d", len(catalog.EventReports))
	log.Printf("   - Vehículos de eventos insertados: %d", len(catalog.EventVehicles))
}

// loadExistingCatalog loads brands and vehicles from database
func loadExistingCatalog(ctx context.Context, pool *pgxpool.Pool, catalog *Catalog) error {
	// Load brands
	rows, err := pool.Query(ctx, "SELECT id, name FROM brands ORDER BY name")
	if err != nil {
		return fmt.Errorf("error querying brands: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b Brand
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return fmt.Errorf("error scanning brand: %w", err)
		}
		catalog.Brands[b.Name] = &b
		catalog.VehiclesByBrandID[b.ID] = []*Vehicle{}
	}

	// Load vehicles
	rows, err = pool.Query(ctx, "SELECT id, brand_id, model_name FROM vehicles ORDER BY model_name")
	if err != nil {
		return fmt.Errorf("error querying vehicles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.BrandID, &v.ModelName); err != nil {
			return fmt.Errorf("error scanning vehicle: %w", err)
		}

		// Find brand name
		brandName := getBrandNameByID(catalog, v.BrandID)
		key := fmt.Sprintf("%s:%s", brandName, v.ModelName)
		catalog.Vehicles[key] = &v
		catalog.VehiclesByBrandID[v.BrandID] = append(catalog.VehiclesByBrandID[v.BrandID], &v)
	}

	return nil
}

// parseVehicleHeaders parses row 1 (brands) and row 2 (vehicle models)
func parseVehicleHeaders(f *excelize.File, catalog *Catalog) error {
	// Row 1: brand names (starting from column T = 20)
	// Row 2: vehicle model names
	brandRow := 1
	modelRow := 2

	// We need to scan from column T (20) to the end
	// Excel columns: T=20, U=21, ..., BH=60, BI=61
	// But vehicle columns end before BH (calificacion) at column 60

	for col := ColVehiculosStart; col < ColCalificacion; col++ {
		cellAddr := columnIndexToName(col) + strconv.Itoa(brandRow)
		brandName, err := f.GetCellValue(SheetName, cellAddr)
		if err != nil {
			continue
		}
		brandName = strings.TrimSpace(brandName)

		cellAddr = columnIndexToName(col) + strconv.Itoa(modelRow)
		modelName, err := f.GetCellValue(SheetName, cellAddr)
		if err != nil {
			continue
		}
		modelName = strings.TrimSpace(modelName)

		// Skip empty cells
		if brandName == "" && modelName == "" {
			continue
		}

		// If brandName is empty, use previous brand (merged cells)
		if brandName == "" {
			// Find previous brand
			for prevCol := col - 1; prevCol >= ColVehiculosStart; prevCol-- {
				if vc, exists := catalog.VehicleColumns[prevCol]; exists {
					brandName = vc.BrandName
					break
				}
			}
		}

		if modelName != "" {
			catalog.VehicleColumns[col] = &VehicleColumn{
				ColIndex:  col,
				BrandName: brandName,
				ModelName: modelName,
			}

			// Ensure brand exists or create new
			if _, exists := catalog.Brands[brandName]; !exists {
				newBrand := &Brand{
					ID:   uuid.New(),
					Name: brandName,
				}
				catalog.Brands[brandName] = newBrand
				catalog.NewBrands = append(catalog.NewBrands, newBrand)
				catalog.VehiclesByBrandID[newBrand.ID] = []*Vehicle{}
			}

			// Ensure vehicle exists or create new
			key := fmt.Sprintf("%s:%s", brandName, modelName)
			if _, exists := catalog.Vehicles[key]; !exists {
				brand := catalog.Brands[brandName]
				newVehicle := &Vehicle{
					ID:        uuid.New(),
					BrandID:   brand.ID,
					ModelName: modelName,
				}
				catalog.Vehicles[key] = newVehicle
				catalog.NewVehicles = append(catalog.NewVehicles, newVehicle)
				catalog.VehiclesByBrandID[brand.ID] = append(catalog.VehiclesByBrandID[brand.ID], newVehicle)
			}
		}
	}

	return nil
}

// parseEventData parses event data rows (starting from row 3)
func parseEventData(f *excelize.File, catalog *Catalog) error {
	rows, err := f.Rows(SheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}
	defer rows.Close()

	rowNum := 0
	for rows.Next() {
		rowNum++
		if rowNum < DataStartRow {
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			log.Printf("⚠️  Error al leer fila %d: %v", rowNum, err)
			continue
		}

		// Skip if row is too short or empty
		if len(row) < ColDistribuidor {
			continue
		}

		// Check if nombre evento is empty (main required field)
		nombreEvento := getCellValue(row, ColNombreEvento)
		if nombreEvento == "" {
			continue
		}

		// Parse event
		event, err := parseEvent(row, catalog)
		if err != nil {
			log.Printf("⚠️  Error al parsear evento en fila %d: %v", rowNum, err)
			continue
		}
		catalog.Events = append(catalog.Events, event)

		// Parse event report
		report, err := parseEventReport(row, event.ID)
		if err != nil {
			log.Printf("⚠️  Error al parsear reporte en fila %d: %v", rowNum, err)
		} else if report != nil {
			catalog.EventReports = append(catalog.EventReports, report)
		}

		// Parse event vehicles
		eventVehicles, err := parseEventVehicles(row, event.ID, catalog)
		if err != nil {
			log.Printf("⚠️  Error al parsear vehículos en fila %d: %v", rowNum, err)
		} else {
			catalog.EventVehicles = append(catalog.EventVehicles, eventVehicles...)
		}
	}

	return nil
}

// parseEvent parses a single event from a row
func parseEvent(row []string, catalog *Catalog) (*Event, error) {
	// Required fields
	marcaEvento := getCellValue(row, ColMarcaEvento)
	nombreEvento := getCellValue(row, ColNombreEvento)

	if marcaEvento == "" || nombreEvento == "" {
		return nil, fmt.Errorf("campos requeridos vacíos")
	}

	// Get brand
	brand, exists := catalog.Brands[marcaEvento]
	if !exists {
		return nil, fmt.Errorf("marca no encontrada: %s", marcaEvento)
	}

	// Parse date
	fechaStr := getCellValue(row, ColFecha)
	startDate, err := parseDate(fechaStr)
	if err != nil {
		log.Printf("⚠️  Error al parsear fecha '%s', usando fecha actual: %v", fechaStr, err)
		startDate = time.Now()
	}

	// Parse year
	anoStr := getCellValue(row, ColAno)
	year := startDate.Year()
	if anoStr != "" {
		if y, err := strconv.Atoi(anoStr); err == nil {
			year = y
		}
	}

	// Parse duration
	diasStr := getCellValue(row, ColDiasEvento)
	dias := 1
	if diasStr != "" {
		if d, err := strconv.Atoi(diasStr); err == nil && d > 0 {
			dias = d
		}
	}

	event := &Event{
		ID:           uuid.New(),
		BrandID:      brand.ID,
		EventType:    getCellValue(row, ColTipoEvento),
		Organizer:    getCellValue(row, ColOrganizador),
		Name:         nombreEvento,
		StartDate:    startDate,
		Year:         year,
		DurationDays: dias,
		State:        getCellValue(row, ColEstado),
		City:         getCellValue(row, ColCiudad),
		Venue:        getCellValue(row, ColSede),
		Dealer:       getCellValue(row, ColDistribuidor),
		Status:       "COMPLETED", // Eventos históricos son COMPLETED
	}

	return event, nil
}

// parseEventReport parses event report data
func parseEventReport(row []string, eventID uuid.UUID) (*EventReport, error) {
	report := &EventReport{
		ID:        uuid.New(),
		EventID:   eventID,
		Completed: true, // Eventos históricos tienen reportes completos
	}

	// Parse optional integer fields
	if val := getCellValue(row, ColEdecanes); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			report.HostessCount = &n
		}
	}

	if val := getCellValue(row, ColAsistencia); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			report.Attendees = &n
		}
	}

	if val := getCellValue(row, ColDinamicas); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			report.ActivitiesCount = &n
		}
	}

	if val := getCellValue(row, ColDatosLevantados); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			report.LeadsCollected = &n
		}
	}

	if val := getCellValue(row, ColProspectos); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			report.Prospects = &n
		}
	}

	if val := getCellValue(row, ColCalificacion); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n >= 1 && n <= 5 {
			report.DealerRating = &n
		}
	}

	// Parse string fields
	if val := getCellValue(row, ColMontaje); val != "" {
		report.SetupVendor = &val
	}

	if val := getCellValue(row, ColComentarios); val != "" {
		report.Comments = &val
	}

	// Parse boolean field (Promocional)
	if val := getCellValue(row, ColPromocional); val != "" {
		hasPromo := parseBool(val)
		report.HasPromotional = &hasPromo
	}

	return report, nil
}

// parseEventVehicles parses vehicle quantities from row
func parseEventVehicles(row []string, eventID uuid.UUID, catalog *Catalog) ([]*EventVehicle, error) {
	var result []*EventVehicle

	for col, vc := range catalog.VehicleColumns {
		if col >= len(row) {
			continue
		}

		qtyStr := getCellValue(row, col)
		if qtyStr == "" {
			continue
		}

		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			continue
		}

		// Get vehicle
		key := fmt.Sprintf("%s:%s", vc.BrandName, vc.ModelName)
		vehicle, exists := catalog.Vehicles[key]
		if !exists {
			log.Printf("⚠️  Vehículo no encontrado: %s", key)
			continue
		}

		eventVehicle := &EventVehicle{
			ID:        uuid.New(),
			EventID:   eventID,
			VehicleID: vehicle.ID,
			Quantity:  qty,
		}
		result = append(result, eventVehicle)
	}

	return result, nil
}

// generateSQLFiles generates SQL files for import
func generateSQLFiles(catalog *Catalog) error {
	// 01_new_brands.sql
	if len(catalog.NewBrands) > 0 {
		if err := writeBrandsSQL(catalog.NewBrands); err != nil {
			return err
		}
	}

	// 02_new_vehicles.sql
	if len(catalog.NewVehicles) > 0 {
		if err := writeVehiclesSQL(catalog.NewVehicles); err != nil {
			return err
		}
	}

	// 03_events.sql
	if err := writeEventsSQL(catalog.Events); err != nil {
		return err
	}

	// 04_event_reports.sql
	if err := writeEventReportsSQL(catalog.EventReports); err != nil {
		return err
	}

	// 05_event_vehicles.sql
	if err := writeEventVehiclesSQL(catalog.EventVehicles); err != nil {
		return err
	}

	return nil
}

// writeBrandsSQL writes brands SQL file
func writeBrandsSQL(brands []*Brand) error {
	path := filepath.Join(OutputDir, "01_new_brands.sql")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("-- New brands from Excel import\n\n")
	for _, b := range brands {
		f.WriteString(fmt.Sprintf("INSERT INTO brands (id, name) VALUES ('%s', '%s');\n",
			b.ID, escapeSQLString(b.Name)))
	}

	return nil
}

// writeVehiclesSQL writes vehicles SQL file
func writeVehiclesSQL(vehicles []*Vehicle) error {
	path := filepath.Join(OutputDir, "02_new_vehicles.sql")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("-- New vehicles from Excel import\n\n")
	for _, v := range vehicles {
		f.WriteString(fmt.Sprintf("INSERT INTO vehicles (id, brand_id, model_name) VALUES ('%s', '%s', '%s');\n",
			v.ID, v.BrandID, escapeSQLString(v.ModelName)))
	}

	return nil
}

// writeEventsSQL writes events SQL file
func writeEventsSQL(events []*Event) error {
	path := filepath.Join(OutputDir, "03_events.sql")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("-- Events from Excel import\n\n")
	for _, e := range events {
		f.WriteString(fmt.Sprintf(
			"INSERT INTO events (id, brand_id, event_type, organizer, name, start_date, year, duration_days, state, city, venue, dealer, status) VALUES "+
				"('%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', '%s', '%s');\n",
			e.ID, e.BrandID,
			escapeSQLString(e.EventType),
			escapeSQLString(e.Organizer),
			escapeSQLString(e.Name),
			e.StartDate.Format("2006-01-02"),
			e.Year, e.DurationDays,
			escapeSQLString(e.State),
			escapeSQLString(e.City),
			escapeSQLString(e.Venue),
			escapeSQLString(e.Dealer),
			e.Status,
		))
	}

	return nil
}

// writeEventReportsSQL writes event reports SQL file
func writeEventReportsSQL(reports []*EventReport) error {
	path := filepath.Join(OutputDir, "04_event_reports.sql")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("-- Event reports from Excel import\n\n")
	for _, r := range reports {
		// Build dynamic SQL based on non-null fields
		cols := []string{"id", "event_id", "completed"}
		vals := []string{
			fmt.Sprintf("'%s'", r.ID),
			fmt.Sprintf("'%s'", r.EventID),
			fmt.Sprintf("%t", r.Completed),
		}

		if r.HostessCount != nil {
			cols = append(cols, "hostess_count")
			vals = append(vals, strconv.Itoa(*r.HostessCount))
		}
		if r.SetupVendor != nil {
			cols = append(cols, "setup_vendor")
			vals = append(vals, fmt.Sprintf("'%s'", escapeSQLString(*r.SetupVendor)))
		}
		if r.HasPromotional != nil {
			cols = append(cols, "has_promotional")
			vals = append(vals, strconv.FormatBool(*r.HasPromotional))
		}
		if r.Attendees != nil {
			cols = append(cols, "attendees")
			vals = append(vals, strconv.Itoa(*r.Attendees))
		}
		if r.ActivitiesCount != nil {
			cols = append(cols, "activities_count")
			vals = append(vals, strconv.Itoa(*r.ActivitiesCount))
		}
		if r.LeadsCollected != nil {
			cols = append(cols, "leads_collected")
			vals = append(vals, strconv.Itoa(*r.LeadsCollected))
		}
		if r.Prospects != nil {
			cols = append(cols, "prospects")
			vals = append(vals, strconv.Itoa(*r.Prospects))
		}
		if r.DealerRating != nil {
			cols = append(cols, "dealer_rating")
			vals = append(vals, strconv.Itoa(*r.DealerRating))
		}
		if r.Comments != nil {
			cols = append(cols, "comments")
			vals = append(vals, fmt.Sprintf("'%s'", escapeSQLString(*r.Comments)))
		}

		f.WriteString(fmt.Sprintf("INSERT INTO event_reports (%s) VALUES (%s);\n",
			strings.Join(cols, ", "),
			strings.Join(vals, ", "),
		))
	}

	return nil
}

// writeEventVehiclesSQL writes event vehicles SQL file
func writeEventVehiclesSQL(eventVehicles []*EventVehicle) error {
	path := filepath.Join(OutputDir, "05_event_vehicles.sql")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("-- Event vehicles from Excel import\n\n")
	for _, ev := range eventVehicles {
		f.WriteString(fmt.Sprintf("INSERT INTO event_vehicles (id, event_id, vehicle_id, quantity) VALUES ('%s', '%s', '%s', %d);\n",
			ev.ID, ev.EventID, ev.VehicleID, ev.Quantity))
	}

	return nil
}

// executeSQL executes SQL in database
func executeSQL(ctx context.Context, pool *pgxpool.Pool, catalog *Catalog) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert new brands
	for _, b := range catalog.NewBrands {
		_, err := tx.Exec(ctx, "INSERT INTO brands (id, name) VALUES ($1, $2)", b.ID, b.Name)
		if err != nil {
			return fmt.Errorf("error inserting brand %s: %w", b.Name, err)
		}
	}

	// Insert new vehicles
	for _, v := range catalog.NewVehicles {
		_, err := tx.Exec(ctx, "INSERT INTO vehicles (id, brand_id, model_name) VALUES ($1, $2, $3)",
			v.ID, v.BrandID, v.ModelName)
		if err != nil {
			return fmt.Errorf("error inserting vehicle %s: %w", v.ModelName, err)
		}
	}

	// Insert events
	for _, e := range catalog.Events {
		_, err := tx.Exec(ctx,
			`INSERT INTO events (id, brand_id, event_type, organizer, name, start_date, year, duration_days, state, city, venue, dealer, status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			e.ID, e.BrandID, e.EventType, e.Organizer, e.Name, e.StartDate, e.Year, e.DurationDays,
			e.State, e.City, e.Venue, e.Dealer, e.Status)
		if err != nil {
			return fmt.Errorf("error inserting event %s: %w", e.Name, err)
		}
	}

	// Insert event reports
	for _, r := range catalog.EventReports {
		_, err := tx.Exec(ctx,
			`INSERT INTO event_reports (id, event_id, hostess_count, setup_vendor, has_promotional, attendees,
			 activities_count, leads_collected, prospects, dealer_rating, comments, completed)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			r.ID, r.EventID, r.HostessCount, r.SetupVendor, r.HasPromotional, r.Attendees,
			r.ActivitiesCount, r.LeadsCollected, r.Prospects, r.DealerRating, r.Comments, r.Completed)
		if err != nil {
			return fmt.Errorf("error inserting event report for event %s: %w", r.EventID, err)
		}
	}

	// Insert event vehicles
	for _, ev := range catalog.EventVehicles {
		_, err := tx.Exec(ctx,
			`INSERT INTO event_vehicles (id, event_id, vehicle_id, quantity) VALUES ($1, $2, $3, $4)`,
			ev.ID, ev.EventID, ev.VehicleID, ev.Quantity)
		if err != nil {
			return fmt.Errorf("error inserting event vehicle: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Utility functions

func getCellValue(row []string, colIndex int) string {
	if colIndex >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[colIndex])
}

func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}

	// Try common formats
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func parseBool(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	return val == "si" || val == "sí" || val == "yes" || val == "true" || val == "1"
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func columnIndexToName(index int) string {
	name := ""
	for index > 0 {
		index--
		name = string(rune('A'+(index%26))) + name
		index /= 26
	}
	return name
}

func getBrandNameByID(catalog *Catalog, brandID uuid.UUID) string {
	for name, brand := range catalog.Brands {
		if brand.ID == brandID {
			return name
		}
	}
	return ""
}
