import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import Chart from 'chart.js/auto';
import { formatDate } from '@angular/common';
import { MatchingStat, MatchingStatService } from '../../services/matching-stat.service';
import { TranslateService } from '@ngx-translate/core';


@Component({
  selector: 'app-statistics',
  templateUrl: './statistics.page.html',
  styleUrls: ['./statistics.page.scss'],
})
export class StatisticsPage implements OnInit {

  @ViewChild('chartCanvas') chartCanvas : ElementRef;
  data : any = [];
  canvasChart : Chart;

  matchingstatistics: MatchingStat[] = [];
  label1: string; label2: string; label3: string;
  numUsers: boolean = false;
  numgroups = 0;
  numconves = 0;
  numfailed = 0;
  numMatchingSection = 0;

  constructor(private matchingstatService: MatchingStatService, public translate: TranslateService) { }

  ngOnInit() {
    this.label1 = this.translate.instant('NUM_CONNEXIONS');
    this.label2 = this.translate.instant('EMAIL_SENT');
    this.label3 = this.translate.instant('EMAIL_FAILED');
    this.getMathingStats();
  }

  getMathingStats() {
    let numgroups = [];
    let numemails = [];
    let numfailed = [];
    let labels = [];

    this.matchingstatService.getMatchingStats().subscribe((matchingstats) => {
      this.numMatchingSection = matchingstats.length;
      this.matchingstatistics = matchingstats
      for(let mathingstat of matchingstats) {
        labels.push(formatDate(mathingstat.CreatedAt, "MMM d, y", "en" /**navigator.language.split("-")[0] ||  */))
        numgroups.push(mathingstat.NumGroups)
        numemails.push(mathingstat.NumConversations)
        numfailed.push(mathingstat.NumFailed)
        this.numconves += mathingstat.NumConversations
        this.numfailed += mathingstat.NumFailed 
      }

      this.data = {
        labels: labels,
        datasets: [{
          label: this.label1,
          data: numgroups,
          backgroundColor: 'rgba(153, 102, 255)',
          borderColor: 'rgb(153, 102, 255)',
          borderWidth: 1
        },
        {
          label: this.label2,
          data: numemails,
          backgroundColor: 'rgba(75, 192, 192)',
          borderColor: 'rgb(75, 192, 192)',
          borderWidth: 1
        },
        {
          label: this.label3,
          data: numfailed,
          backgroundColor: 'rgba(255, 99, 132)',
          borderColor: 'rgb(255, 99, 132)',
          borderWidth: 1
        }]
      };
      this.changeChart({detail: {
        value : 'bar'
      }});
    });
  }

  changeChart( event: any ) {
    const type = event.detail.value || 'bar';
    if ( this.canvasChart ) {
      this.canvasChart.destroy();
    }
    this.canvasChart = new Chart(this.chartCanvas.nativeElement, {
      type: type,
      data: this.data,
      options: {
        indexAxis: 'x'
      }
    });
  }

  addRandom( points: any ) : number {
    return Number(points) - Number( Math.floor((Math.random() * 100) + 1) );
  }
}
