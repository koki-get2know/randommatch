import { TestBed } from '@angular/core/testing';

import { MatchingStatService } from './matching-stat.service';

describe('MatchingStatService', () => {
  let service: MatchingStatService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MatchingStatService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
