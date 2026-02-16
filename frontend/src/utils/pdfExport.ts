import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';
import { Event, DashboardAnalytics } from '../types';

/**
 * Exporta un evento individual a PDF
 */
export const exportEventToPDF = async (event: Event, elementId: string = 'event-detail-content') => {
  try {
    // Capturar el elemento HTML
    const element = document.getElementById(elementId);
    if (!element) {
      throw new Error('Elemento no encontrado para exportar');
    }

    // Convertir a canvas
    const canvas = await html2canvas(element, {
      scale: 2,
      logging: false,
      useCORS: true,
    });

    // Crear PDF
    const imgWidth = 190; // A4 width in mm minus margins
    const pageHeight = 277; // A4 height in mm
    const imgHeight = (canvas.height * imgWidth) / canvas.width;

    const pdf = new jsPDF('p', 'mm', 'a4');
    const imgData = canvas.toDataURL('image/png');

    let heightLeft = imgHeight;
    let position = 10;

    // Agregar imagen al PDF
    pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
    heightLeft -= pageHeight;

    // Agregar páginas adicionales si es necesario
    while (heightLeft > 0) {
      position = heightLeft - imgHeight + 10;
      pdf.addPage();
      pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
      heightLeft -= pageHeight;
    }

    // Generar nombre de archivo
    const fileName = `evento_${event.name.replace(/\s+/g, '_')}_${new Date().toISOString().split('T')[0]}.pdf`;

    // Descargar
    pdf.save(fileName);

    return true;
  } catch (error) {
    console.error('Error al exportar evento a PDF:', error);
    throw error;
  }
};

/**
 * Exporta el dashboard de analytics a PDF
 */
export const exportDashboardToPDF = async (
  analytics: DashboardAnalytics,
  filters?: { year?: number; brandId?: string }
) => {
  try {
    const pdf = new jsPDF('p', 'mm', 'a4');
    const pageWidth = 210;
    const pageHeight = 297;
    const margin = 15;
    let currentY = margin;

    // Función helper para agregar título
    const addTitle = (text: string, fontSize: number = 16) => {
      pdf.setFontSize(fontSize);
      pdf.setFont('helvetica', 'bold');
      pdf.text(text, margin, currentY);
      currentY += 8;
    };

    // Función helper para agregar texto normal
    const addText = (text: string, fontSize: number = 10) => {
      pdf.setFontSize(fontSize);
      pdf.setFont('helvetica', 'normal');
      pdf.text(text, margin, currentY);
      currentY += 6;
    };

    // Función helper para verificar nueva página
    const checkNewPage = (spaceNeeded: number = 20) => {
      if (currentY + spaceNeeded > pageHeight - margin) {
        pdf.addPage();
        currentY = margin;
      }
    };

    // Título principal
    addTitle('Reporte de Analytics - CustomerMX', 18);

    // Fecha y filtros
    pdf.setFontSize(9);
    pdf.setTextColor(100);
    addText(`Generado: ${new Date().toLocaleDateString('es-MX', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })}`);

    if (filters?.year) {
      addText(`Año: ${filters.year}`);
    }

    currentY += 5;
    pdf.setTextColor(0);

    // KPIs principales
    checkNewPage(40);
    addTitle('Métricas Generales', 14);

    const kpis = [
      { label: 'Total de Eventos', value: analytics.totals.total_events },
      { label: 'Total de Asistentes', value: analytics.totals.total_attendees.toLocaleString() },
      { label: 'Total de Leads', value: analytics.totals.total_leads.toLocaleString() },
      { label: 'Total de Prospectos', value: analytics.totals.total_prospects.toLocaleString() },
      { label: 'Promedio de Asistentes', value: analytics.totals.average_attendees.toFixed(0) },
      { label: 'Rating Promedio', value: analytics.totals.average_rating.toFixed(2) },
    ];

    kpis.forEach(kpi => {
      pdf.setFont('helvetica', 'bold');
      pdf.text(`${kpi.label}:`, margin, currentY);
      pdf.setFont('helvetica', 'normal');
      pdf.text(String(kpi.value), margin + 70, currentY);
      currentY += 6;
    });

    // Métricas por marca
    if (analytics.by_brand.length > 0) {
      checkNewPage(50);
      currentY += 5;
      addTitle('Métricas por Marca', 14);

      analytics.by_brand.forEach((brand, index) => {
        checkNewPage(15);
        pdf.setFont('helvetica', 'bold');
        addText(`${index + 1}. ${brand.brand_name}`);
        pdf.setFont('helvetica', 'normal');
        addText(`   Eventos: ${brand.event_count} | Asistentes: ${brand.total_attendees.toLocaleString()} | Rating: ${brand.average_rating.toFixed(2)}`);
      });
    }

    // Top vehículos
    if (analytics.top_vehicles.length > 0) {
      checkNewPage(50);
      currentY += 5;
      addTitle('Top 10 Vehículos Presentados', 14);

      analytics.top_vehicles.slice(0, 10).forEach((vehicle, index) => {
        checkNewPage(15);
        pdf.setFont('helvetica', 'bold');
        addText(`${index + 1}. ${vehicle.model_name} (${vehicle.brand_name})`);
        pdf.setFont('helvetica', 'normal');
        addText(`   Presentaciones: ${vehicle.times_presented} | Cantidad total: ${vehicle.total_quantity}`);
      });
    }

    // Métricas de conversión
    checkNewPage(40);
    currentY += 5;
    addTitle('Métricas de Conversión', 14);

    const conversion = analytics.conversion;
    addText(`Asistentes Totales: ${conversion.total_attendees.toLocaleString()}`);
    addText(`Leads Capturados: ${conversion.total_leads.toLocaleString()} (${conversion.lead_rate.toFixed(1)}%)`);
    addText(`Prospectos Calificados: ${conversion.total_prospects.toLocaleString()} (${conversion.prospect_rate.toFixed(1)}%)`);

    // Top distribuidores
    if (analytics.top_dealers.length > 0) {
      checkNewPage(50);
      currentY += 5;
      addTitle('Top Distribuidores', 14);

      analytics.top_dealers.slice(0, 10).forEach((dealer, index) => {
        checkNewPage(15);
        pdf.setFont('helvetica', 'bold');
        addText(`${index + 1}. ${dealer.dealer}`);
        pdf.setFont('helvetica', 'normal');
        addText(`   Rating: ${dealer.average_rating.toFixed(2)} | Eventos: ${dealer.event_count} | Asistentes: ${dealer.total_attendees.toLocaleString()}`);
      });
    }

    // Distribución por estado
    if (analytics.by_state.length > 0) {
      checkNewPage(50);
      currentY += 5;
      addTitle('Distribución por Estado', 14);

      analytics.by_state.forEach((state, index) => {
        checkNewPage(10);
        addText(`${state.state}: ${state.event_count} eventos, ${state.attendees.toLocaleString()} asistentes`);
      });
    }

    // Tipos de eventos
    if (analytics.by_event_type.length > 0) {
      checkNewPage(40);
      currentY += 5;
      addTitle('Distribución por Tipo de Evento', 14);

      analytics.by_event_type.forEach((type) => {
        checkNewPage(10);
        addText(`${type.event_type}: ${type.event_count} eventos, ${type.attendees.toLocaleString()} asistentes`);
      });
    }

    // Pie de página en cada página
    const pageCount = pdf.getNumberOfPages();
    for (let i = 1; i <= pageCount; i++) {
      pdf.setPage(i);
      pdf.setFontSize(8);
      pdf.setTextColor(150);
      pdf.text(
        `Página ${i} de ${pageCount}`,
        pageWidth / 2,
        pageHeight - 10,
        { align: 'center' }
      );
      pdf.text(
        'CustomerMX - Sistema de Gestión de Eventos',
        pageWidth - margin,
        pageHeight - 10,
        { align: 'right' }
      );
    }

    // Generar nombre de archivo
    const fileName = `analytics_dashboard_${filters?.year || 'todos'}_${new Date().toISOString().split('T')[0]}.pdf`;

    // Descargar
    pdf.save(fileName);

    return true;
  } catch (error) {
    console.error('Error al exportar dashboard a PDF:', error);
    throw error;
  }
};

/**
 * Captura un elemento específico y lo convierte a imagen para PDF
 */
export const captureElementAsImage = async (elementId: string): Promise<string> => {
  const element = document.getElementById(elementId);
  if (!element) {
    throw new Error(`Elemento con ID "${elementId}" no encontrado`);
  }

  const canvas = await html2canvas(element, {
    scale: 2,
    logging: false,
    useCORS: true,
    backgroundColor: '#ffffff',
  });

  return canvas.toDataURL('image/png');
};
