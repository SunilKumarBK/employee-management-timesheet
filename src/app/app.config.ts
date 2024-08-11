import { ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideClientHydration } from '@angular/platform-browser';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideHttpClient, withFetch, withInterceptors } from '@angular/common/http';
import { spinnerInterceptor } from './loader/interceptors/spinner.interceptor';
import { BrowserAnimationsModule, provideAnimations } from '@angular/platform-browser/animations';
import { provideToastr } from 'ngx-toastr';
import { NgxOrgChartModule } from 'ngx-org-chart';

export const appConfig: ApplicationConfig = {
  providers: [provideZoneChangeDetection({ eventCoalescing: true }), provideRouter(routes), provideClientHydration(), provideAnimationsAsync(),provideHttpClient(
    withFetch()),provideHttpClient(withInterceptors([spinnerInterceptor]))  , provideAnimationsAsync(),importProvidersFrom([BrowserAnimationsModule]),provideToastr(),provideAnimations(),NgxOrgChartModule]
};
