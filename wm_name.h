#ifndef WM_NAME_H
#define WM_NAME_H

#include <xcb/xcb.h>

typedef struct {
	xcb_connection_t *conn;
	xcb_screen_t     *screen;
	int               screen_num;
} xconn_t;

xconn_t *connect_x(void);
void disconnect_x(xconn_t *xconn);
int set_screen(xconn_t *xconn);
int set_wm_name(xconn_t *xconn, const char *name);

#endif
