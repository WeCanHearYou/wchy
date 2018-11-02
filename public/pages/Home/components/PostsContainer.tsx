import React from "react";

import { ListPosts, TagsFilter, PostFilter } from "../";

import { Post, Tag, CurrentUser } from "@fider/models";
import { Loader, Field, Input } from "@fider/components";
import { actions, navigator, querystring } from "@fider/services";

interface PostsContainerProps {
  user?: CurrentUser;
  posts: Post[];
  tags: Tag[];
  countPerStatus: { [key: string]: number };
}

interface PostsContainerState {
  loading: boolean;
  posts?: Post[];
  view: string;
  tags: string[];
  query: string;
  limit?: number;
}

export class PostsContainer extends React.Component<PostsContainerProps, PostsContainerState> {
  constructor(props: PostsContainerProps) {
    super(props);

    this.state = {
      posts: this.props.posts,
      loading: false,
      view: querystring.get("view"),
      query: querystring.get("query"),
      tags: querystring.getArray("tags"),
      limit: querystring.getNumber("limit")
    };
  }

  private changeFilterCriteria<K extends keyof PostsContainerState>(
    obj: Pick<PostsContainerState, K>,
    reset: boolean
  ): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase();
      navigator.replaceState(
        querystring.stringify({
          tags: this.state.tags,
          query,
          view: this.state.view,
          limit: this.state.limit
        })
      );

      this.searchPosts(query, this.state.view, this.state.limit, this.state.tags, reset);
    });
  }

  private timer?: number;
  private async searchPosts(query: string, view: string, limit: number | undefined, tags: string[], reset: boolean) {
    window.clearTimeout(this.timer);
    this.setState({ posts: reset ? undefined : this.state.posts, loading: true });
    this.timer = window.setTimeout(() => {
      actions.searchPosts({ query, view, limit, tags }).then(response => {
        if (this.state.loading) {
          this.setState({ loading: false, posts: response.data });
        }
      });
    }, 200);
  }

  private handleViewChanged = (view: string) => {
    this.changeFilterCriteria({ view }, true);
  };

  private handleTagsFilterChanged = (tags: string[]) => {
    this.changeFilterCriteria({ tags }, true);
  };

  private handleSearchFilterChanged = (query: string) => {
    this.changeFilterCriteria({ query }, true);
  };

  private handleSearchClick = (query: string) => {
    this.changeFilterCriteria({ query }, true);
  };

  private clearSearch = () => {
    this.changeFilterCriteria({ query: "" }, true);
  };

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault();
    this.changeFilterCriteria({ limit: (this.state.limit || 30) + 10 }, false);
  };

  private canShowMore = (): boolean => {
    return this.state.posts ? this.state.posts.length >= (this.state.limit || 30) : false;
  };

  public render() {
    return (
      <>
        <div className="row">
          {!this.state.query && (
            <div className="l-filter-col col-sm-7 col-md-8 col-lg-9 mb-2">
              <Field>
                <PostFilter
                  activeView={this.state.view}
                  viewChanged={this.handleViewChanged}
                  countPerStatus={this.props.countPerStatus}
                />
                <TagsFilter
                  tags={this.props.tags}
                  selectionChanged={this.handleTagsFilterChanged}
                  defaultSelection={this.state.tags}
                />
              </Field>
            </div>
          )}
          <div className={!this.state.query ? `l-search-col col-sm-5 col-md-4 col-lg-3 mb-2` : "col-sm-12 mb-2"}>
            <Input
              field="query"
              icon={this.state.query ? "cancel" : "search"}
              onIconClick={this.state.query ? this.clearSearch : undefined}
              placeholder="Search..."
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            />
          </div>
        </div>
        <ListPosts
          posts={this.state.posts}
          tags={this.props.tags}
          emptyText={"No results matched your search, try something different."}
        />
        {this.state.loading && <Loader />}
        {this.canShowMore() && (
          <h5 className="c-post-list-show-more" onTouchEnd={this.showMore} onClick={this.showMore}>
            View more posts
          </h5>
        )}
      </>
    );
  }
}
