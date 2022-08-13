import { TestBed } from '@angular/core/testing';

import { BearerInterceptor } from './bearer-interceptor.service';

describe('BearerInterceptorService', () => {
  let service: BearerInterceptor;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BearerInterceptor);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
